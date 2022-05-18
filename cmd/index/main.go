package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-properties"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/sfomuseum/go-flags/multi"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {

	props := flag.String("properties", "", "The path to your whosonfirst-properties/properties directory")
	iterator_uri := flag.String("iterator-uri", "repo://", "A valid go-whosonfirst-iterate/v2 URI.")
	debug := flag.Bool("debug", false, "Go through all the motions but don't write any new files")

	var alternate multi.MultiString
	flag.Var(&alternate, "alternate", "One or more paths to alternate properties directories")
	
	flag.Parse()

	ctx := context.Background()

	mu := new(sync.RWMutex)
	seen := new(sync.Map)

	if len(alternate) > 0 {

		alternate_cb := func(path string, info os.FileInfo) error {

			if info.IsDir() {
				return nil
			}
			
			if filepath.Ext(path) != ".json" {
				return nil
			}
			
			prop, err := properties.NewPropertyFromFile(path)

			if err != nil {
				return err
			}

			seen.Store(prop.String(), true)
			return nil
		}

		for _, path := range alternate {

			cr := crawl.NewCrawler(path)
			err := cr.Crawl(alternate_cb)

			if err != nil {
				log.Fatalf("Failed to crawl alternate properties source '%s', %v", path, err)
			}
			
		}
	}
	
	cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		select {

		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		body, err := io.ReadAll(fh)

		if err != nil {
			return fmt.Errorf("Unable to load %s, because %s", path, err)
		}

		pr := gjson.GetBytes(body, "properties")

		if !pr.Exists() {
			return fmt.Errorf("%s is missing a properties dictionary!", path)
		}

		// PLEASE FOR TO go func() ME...

		for k, _ := range pr.Map() {

			_, exists := seen.Load(k)

			if exists {
				continue
			}

			p, err := properties.NewPropertyFromKey(k)

			if err != nil {
				log.Printf("failed to parse key (%s) for %s\n", k, path)
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

						return fmt.Errorf("failed to ensure ID for %s, because %v", abs_path, err)
					}

					err = p.Write(*props)

					if err != nil {
						log.Printf("failed to write (%s) for %s, because %v\n", abs_path, path, err)
					}
				}
			}

			seen.Store(k, true)
			mu.Unlock()
		}

		return nil

	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, cb)

	if err != nil {
		log.Fatal(err)
	}

	uris := flag.Args()

	err = iter.IterateURIs(ctx, uris...)

	if err != nil {
		log.Fatal(err)
	}
}
