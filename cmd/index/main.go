package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-properties"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {

	props := flag.String("properties", "", "The path to your whosonfirst-properties/properties directory")
	iterator_uri := flag.String("iterator-uri", "repo://", "A valid go-whosonfirst-iterate/emitter URI.")
	debug := flag.Bool("debug", false, "Go through all the motions but don't write any new files")

	flag.Parse()

	ctx := context.Background()

	mu := new(sync.RWMutex)
	seen := make(map[string]bool)

	cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		select {

		case <-ctx.Done():
			return nil
		default:
			path, err := emitter.PathForContext(ctx)

			if err != nil {
				return err
			}

			_, uri_args, err := uri.ParseURI(path)

			if err != nil {
				return err
			}

			if uri_args.IsAlternate {
				return nil
			}

			closer := io.NopCloser(fh)

			f, err := feature.LoadWOFFeatureFromReader(closer)

			if err != nil {
				msg := fmt.Sprintf("Unable to load %s, because %s", path, err)
				return errors.New(msg)
			}

			pr := gjson.GetBytes(f.Bytes(), "properties")

			if !pr.Exists() {
				msg := fmt.Sprintf("%s is missing a properties dictionary!", f.Id())
				return errors.New(msg)
			}

			// PLEASE FOR TO go func() ME...

			for k, _ := range pr.Map() {

				mu.RLock()
				_, ok := seen[k]
				mu.RUnlock()

				if ok {
					continue
				}

				p, err := properties.NewPropertyFromKey(k)

				if err != nil {
					msg := fmt.Sprintf("failed to parse key (%s) for %s", k, f.Id())
					log.Println(msg)
					continue
				}

				if p.IsName() {

					if *debug {
						log.Printf("%s is a name property, skipping\n", p)
					}

					continue
				}

				rel_path := p.RelPath()
				abs_path := filepath.Join(*props, rel_path)

				mu.Lock()

				_, err = os.Stat(abs_path)

				if os.IsNotExist(err) {

					if *debug {
						log.Printf("create %s but debugging is enabled, so don't\n", abs_path)
					} else {
						err = p.EnsureId()

						if err != nil {

							mu.Unlock()

							msg := fmt.Sprintf("failed to ensure ID for %s, because %s", abs_path, err)
							return errors.New(msg)
						}

						err = p.Write(*props)

						if err != nil {
							msg := fmt.Sprintf("failed to write (%s) for %s, because", abs_path, f.Id(), err)
							log.Println(msg)
						}
					}
				}

				seen[k] = true
				mu.Unlock()
			}

			return nil
		}

	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, cb)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	err = iter.IterateURIs(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}
}
