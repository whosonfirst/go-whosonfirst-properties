package index

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/whosonfirst/go-whosonfirst-properties"	
)

// CatalogPropertiesOptions is a struct containing configuration data for the `CatalogProperties` method.
type CatalogPropertiesOptions struct {
	// Lookup is a `sync.Map` instance whose keys are the names of properties that have already been encountered (crawl property definition files).
	Lookup *sync.Map
}

// CatalogProperties() will crawl one or more directories containing Who's On First style property definition
// files and cataloging each match in a `sync.Map` instance.
func CatalogProperties(ctx context.Context, opts *CatalogPropertiesOptions, paths ...string) error {

	for _, path := range paths {
		
		dir_fs := os.DirFS(path)

		fs.WalkDir(dir_fs, ".", func(path string, info fs.DirEntry, err error) error {
			
			if err != nil {
				log.Fatal(err)
			}
			
			if info.IsDir() {
				return nil
			}
			
			if filepath.Ext(path) != ".json" {
				return nil
			}
			
			prop, err := properties.NewPropertyFromFile(path)
			
			if err != nil {
				slog.Warn("Failed to parse as properties file", "path", path, "error", err)
				return nil
			}

			opts.Lookup.Store(prop.String(), true)
			return nil
		})
	}

	return nil
}
