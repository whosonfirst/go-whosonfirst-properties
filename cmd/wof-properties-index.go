package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-index/utils"
	"github.com/whosonfirst/go-whosonfirst-properties"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Closer struct {
	fh io.Reader
}

func (c Closer) Read(b []byte) (int, error) {
	return c.fh.Read(b)
}

func (c Closer) Close() error {
	return nil
}

func main() {

	valid_modes := strings.Join(index.Modes(), ",")
	desc_modes := fmt.Sprintf("The mode to use importing data. Valid modes are: %s.", valid_modes)

	props := flag.String("properties", "", "The path to your whosonfirst-properties/properties directory")
	mode := flag.String("mode", "repo", desc_modes)
	debug := flag.Bool("debug", false, "Go through all the motions but don't write any new files")

	flag.Parse()

	mu := new(sync.RWMutex)
	seen := make(map[string]bool)

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		select {

		case <-ctx.Done():
			return nil
		default:
			path, err := index.PathForContext(ctx)

			if err != nil {
				return err
			}

			ok, err := utils.IsPrincipalWOFRecord(fh, ctx)

			if err != nil {
				return err
			}

			if !ok {
				return nil
			}

			// HACK - see above
			closer := Closer{fh}

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

	i, err := index.NewIndexer(*mode, cb)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	err = i.IndexPaths(paths)

	if err != nil {
		log.Fatal(err)
	}
}
