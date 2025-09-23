//go:build wasmjs
package wasm

import (
	"fmt"
	"syscall/js"
	"sync"
	"strings"
	
	"github.com/sfomuseum/go-csvdict/v2"
)

func LookupFunc() js.Func {

	lookup_func := sync.OnceValues(func() (*sync.Map, error) {

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
	
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		
		prop_k := args[0].String()
		prop_k = strings.ToLower(prop_k)
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			
			resolve := args[0]
			reject := args[1]
			
			lookup, err := lookup_func()
			
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to build lookup, %w", err))
				return nil
			}
			
			_, exists := lookup.Load(prop_k)
			
			if !exists {
				reject.Invoke()
				return nil
			}
			
			resolve.Invoke()
			return nil
		})
		
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
