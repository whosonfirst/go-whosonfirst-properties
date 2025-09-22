// index is a command line tool for crawling one or more Who's On First data sources and ensuring that
// individual properties contained in those records have a corresponding machine-readable properties
// description file.
package main

import ()

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/whosonfirst/go-whosonfirst-iterate-git/v3/github"

	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	"github.com/whosonfirst/go-whosonfirst-properties"
	"github.com/whosonfirst/go-whosonfirst-properties/index"
	"log"
	"sync"
)

func main() {

	var root string
	var debug bool

	iterator_uri := flag.String("iterator-uri", "repo://", "A valid go-whosonfirst-iterate/v2 URI.")

	flag.StringVar(&root, "properties", "", "The path to your whosonfirst-properties/properties directory")
	flag.BoolVar(&debug, "debug", false, "Go through all the motions but don't write any new files.")

	var alternates multi.MultiString
	flag.Var(&alternates, "alternate", "One or more paths to alternate properties directories that will be crawled to check for existing properties (that will not be duplicated).")

	var exclude multi.MultiRegexp
	flag.Var(&exclude, "exclude", "One or more valid regular expressions to use for excluding property names you don't want to index")

	flag.Parse()

	ctx := context.Background()

	lookup := new(sync.Map)

	if len(alternates) > 0 {

		crawl_alternates_opts := &index.CatalogPropertiesOptions{
			Lookup: lookup,
		}

		err := index.CatalogProperties(ctx, crawl_alternates_opts, alternates...)

		if err != nil {
			log.Fatalf("Failed to crawl alternate sources, %v", err)
		}
	}

	iter, err := iterate.NewIterator(ctx, *iterator_uri)

	if err != nil {
		log.Fatalf("Failed to create new iterator, %v", err)
	}

	uris := flag.Args()

	for rec, err := range iter.Iterate(ctx, uris...) {

		if err != nil {
			log.Fatal(err)
		}

		defer rec.Body.Close()

		body, err := io.ReadAll(rec.Body)

		if err != nil {
			log.Fatalf("Unable to load %s, because %s", rec.Path, err)
		}

		pr := gjson.GetBytes(body, "properties")

		if !pr.Exists() {
			log.Fatalf("%s is missing a properties dictionary!", rec.Path)
		}

		// PLEASE FOR TO go func() ME...

		for k, _ := range pr.Map() {

			_, exists := lookup.Load(k)

			if exists {
				continue
			}

			p, err := properties.NewPropertyFromKey(k)

			if err != nil {
				slog.Warn("Failed to parse key", "path", rec.Path, "key", k, "error", err)
				continue
			}

			if p.IsName() {

				if debug {
					slog.Debug("Duplicate, skipping")
				}

				continue
			}

			if len(exclude) > 0 {

				include := true

				for _, re := range exclude {

					if re.MatchString(p.String()) {
						include = false
						break
					}
				}

				if !include {
					continue
				}
			}

			// START OF should be updated to use gocloud.dev/blob
			// and sfomuseum/go-atomicwrite - or possibly whosonfirst/go-writer
			// which would allow writes directly back to github but it's also
			// not clear that's really necessary or desired...

			rel_path := p.RelPath()
			abs_path := filepath.Join(root, rel_path)

			_, err = os.Stat(abs_path)

			if os.IsNotExist(err) {

				if debug {
					slog.Debug("create path but debugging is enabled, so don't", "path", abs_path)
				} else {

					err = p.EnsureId()

					if err != nil {
						log.Fatalf("failed to ensure ID for %s, because %v", abs_path, err)
					}

					err = p.Write(root)

					if err != nil {
						slog.Warn("Failed to write property", "path", abs_path, "error", err)
						continue
					}
				}
			}

			// END OF should be updated to use gocloud.dev/blob

			lookup.Store(k, true)
		}

	}
}
