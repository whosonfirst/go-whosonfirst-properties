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

type CrawlAlternatesOptions struct {
	Lookup *sync.Map
	Logger *log.Logger
}

func CrawlAlternates(ctx context.Context, opts *CrawlAlternatesOptions, alternates ...string) error {

	alternate_cb := func(path string, info os.FileInfo) error {

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

	for _, path := range alternates {

		cr := crawl.NewCrawler(path)
		err := cr.Crawl(alternate_cb)

		if err != nil {
			return fmt.Errorf("Failed to crawl alternate properties source '%s', %v", path, err)
		}

	}

	return nil
}
