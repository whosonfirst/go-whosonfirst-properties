package main

import (
	"flag"
	"io"
	"io/fs"
	"log"
	_ "log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sfomuseum/go-csvdict/v2"
	"github.com/whosonfirst/go-whosonfirst-properties"
)

func main() {

	var root string
	var report string

	flag.StringVar(&root, "properties", "", "The path to your whosonfirst-properties/properties directory")
	flag.StringVar(&report, "report", "", "The path to write your whosonfirst-properties report. Default is STDOUT.")

	flag.Parse()

	_, err := os.Stat(root)

	if err != nil {
		log.Fatal(err)
	}

	var wr io.WriteCloser

	if report == "" {
		wr = os.Stdout
	} else {
		f_wr, err := os.OpenFile(report, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatal(err)
		}

		wr = f_wr
	}

	csv_wr, err := csvdict.NewWriter(wr)

	if err != nil {
		log.Fatal(err)
	}

	props_fs := os.DirFS(root)

	err = fs.WalkDir(props_fs, ".", func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		r, err := props_fs.Open(path)

		if err != nil {
			return err
		}

		defer r.Close()

		prop, err := properties.NewPropertyFromReader(r)

		if err != nil {
			return err
		}

		row := map[string]string{
			"id":          strconv.FormatInt(prop.Id, 10),
			"prefix":      prop.Prefix,
			"name":        prop.Name,
			"description": prop.Description,
		}

		err = csv_wr.WriteRow(row)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	csv_wr.Flush()

	if report == "" {
		err := wr.Close()

		if err != nil {
			log.Fatal(err)
		}
	}
}
