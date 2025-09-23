//go:build wasmjs
package wasm

import (
	"fmt"
	"syscall/js"
	"strings"
)

func LookupFunc() js.Func {
	
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		
		prop_k := args[0].String()
		prop_k = strings.ToLower(prop_k)
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			
			resolve := args[0]
			reject := args[1]
			
			lookup, err := properties_func()
			
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
