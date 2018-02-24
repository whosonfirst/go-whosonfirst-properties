package main

import (
	"encoding/csv"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-properties"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {

	props := flag.String("properties", "", "The path to your whosonfirst-properties/properties directory")
	report := flag.String("report", "", "The path to write your whosonfirst-properties report. Default is STDOUT.")

	flag.Parse()

	_, err := os.Stat(*props)

	if err != nil {
		log.Fatal(err)
	}

	var fh io.Writer

	if *report == "" {
		fh = os.Stdout
	} else {
		f, err := os.OpenFile(*report, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		fh = f
	}

	wr := csv.NewWriter(fh)
	mu := new(sync.Mutex)

	row := []string{
		"id",
		"prefix",
		"name",
		"description",
	}

	wr.Write(row)

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

		mu.Lock()
		defer mu.Unlock()

		row := []string{
			strconv.FormatInt(prop.Id, 10),
			prop.Prefix,
			prop.Name,
			prop.Description,
		}

		err = wr.Write(row)

		if err != nil {
			return err
		}

		return nil
	}

	cr := crawl.NewCrawler(*props)
	err = cr.Crawl(cb)

	if err != nil {
		log.Fatal(err)
	}

	wr.Flush()

	err = wr.Error()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
