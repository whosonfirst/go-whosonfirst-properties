package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-properties"
	"log"
	"os"
	"path/filepath"
)

func main() {

	props := flag.String("properties", "", "")

	flag.Parse()

	cb := func(path string, info os.FileInfo) error {

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

		log.Println(path, prop)
		return nil
	}

	cr := crawl.NewCrawler(*props)
	err := cr.Crawl(cb)

	if err != nil {
		log.Fatal(err)
	}
}
