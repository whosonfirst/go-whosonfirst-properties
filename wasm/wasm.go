//go:build wasmjs
package wasm

import (
	"embed"
)

//go:embed props.csv
var FS embed.FS
