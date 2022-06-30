package index

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-properties"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// CatalogPropertiesOptions is a struct containing configuration data for the `CatalogProperties` method.
type CatalogPropertiesOptions struct {
	// Lookup is a `sync.Map` instance whose keys are the names of properties that have already been encountered (crawl property definition files).
	Lookup *sync.Map
	// Logger is a `log.Logger` instance used to log state and feedback.
	Logger *log.Logger
}

// CatalogProperties() will crawl one or more directories containing Who's On First style property definition
// files and cataloging each match in a `sync.Map` instance.
func CatalogProperties(ctx context.Context, opts *CatalogPropertiesOptions, paths ...string) error {

	cb := func(path string, info os.FileInfo) error {

		select {

		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		prop, err := properties.NewPropertyFromFile(path)

		if err != nil {
			opts.Logger.Printf("Failed to parse '%s' as a properties file, %v", path, err)
			return nil
		}

		opts.Lookup.Store(prop.String(), true)
		return nil
	}

	for _, path := range paths {

		cr := crawl.NewCrawler(path)
		err := cr.Crawl(cb)

		if err != nil {
			return fmt.Errorf("Failed to crawl properties source '%s', %v", path, err)
		}

	}

	return nil
}
