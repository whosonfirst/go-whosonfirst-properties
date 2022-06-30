// index is a command line tool for crawling one or more Who's On First data sources and ensuring that
// individual properties contained in those records have a corresponding machine-readable properties
// description file.
package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-iterate-organization"
)

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-properties/index"
	"log"
	"sync"
)

func main() {

	props := flag.String("properties", "", "The path to your whosonfirst-properties/properties directory")
	iterator_uri := flag.String("iterator-uri", "repo://", "A valid go-whosonfirst-iterate/v2 URI.")
	debug := flag.Bool("debug", false, "Go through all the motions but don't write any new files.")

	var alternates multi.MultiString
	flag.Var(&alternates, "alternate", "One or more paths to alternate properties directories that will be crawled to check for existing properties (that will not be duplicated).")

	var exclude multi.MultiRegexp
	flag.Var(&exclude, "exclude", "One or more valid regular expressions to use for excluding property names you don't want to index")

	flag.Parse()

	ctx := context.Background()

	mu := new(sync.RWMutex)
	lookup := new(sync.Map)

	logger := log.Default()

	if len(alternates) > 0 {

		crawl_alternates_opts := &index.CatalogPropertiesOptions{
			Lookup: lookup,
			Logger: logger,
		}

		err := index.CatalogProperties(ctx, crawl_alternates_opts, alternates...)

		if err != nil {
			log.Fatalf("Failed to crawl alternate sources, %v", err)
		}
	}

	iter_cb_opts := &index.EmitterCallbackFuncOptions{
		Debug:   *debug,
		Lookup:  lookup,
		Mutex:   mu,
		Root:    *props,
		Exclude: exclude,
		Logger:  logger,
	}

	iter_cb := index.EmitterCallbackFunc(iter_cb_opts)

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create new iterator, %v", err)
	}

	uris := flag.Args()

	err = iter.IterateURIs(ctx, uris...)

	if err != nil {
		log.Fatalf("Failed to iterate URIs, %v", err)
	}
}
