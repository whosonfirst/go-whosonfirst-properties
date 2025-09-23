//go:build wasmjs
package main

import (
	"log"
	"syscall/js"

	"github.com/whosonfirst/go-whosonfirst-properties/wasm"	
)

func main() {

	lookup_func := wasm.LookupFunc()

	defer lookup_func.Release()

	js.Global().Set("wof_properties_lookup", lookup_func)

	c := make(chan struct{}, 0)

	log.Println("Who's On First properties WASM binary initialized")
	<-c
}
