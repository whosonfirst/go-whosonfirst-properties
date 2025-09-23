//go:build wasmjs
package wasm

import (
	"embed"
	"sync"
	"fmt"
	"strings"
	
	"github.com/sfomuseum/go-csvdict/v2"	
)

//go:embed props.csv
var FS embed.FS

var properties_func = sync.OnceValues(func() (*sync.Map, error) {
	
	r, err := FS.Open("props.csv")
	
	if err != nil {
		return nil, err
	}
	
	defer r.Close()
	
	csv_r, err := csvdict.NewReader(r)
	
	if err != nil {
		return nil, err
	}
	
	lookup := new(sync.Map)
	wg := new(sync.WaitGroup)
	
	for row, err := range csv_r.Iterate(){
		
		if err != nil {
			return nil, err
		}
		
		wg.Go(func(){
			k := fmt.Sprintf("%s:%s", row["prefix"], row["name"])
			k = strings.ToLower(k)
			lookup.Store(k, true)
		})
	}
	
	wg.Wait()
	return lookup, nil
})
