package index

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-properties"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// EmitterCallbackFuncOptions is a struct containing configuration options for the `EmitterCallbackFunc` method.
type EmitterCallbackFuncOptions struct {
	// Debug is a boolean flag to signal that records should not be created or updated.
	Debug bool
	// Lookup is a `sync.Map` instance whose keys are the names of properties that have already been seen or processed.
	Lookup *sync.Map
	// Root is the root directory (path) where new and updated properties should be written.
	Root string
	// Mutex us a `sync.RWMutex` instance used prevent duplicate processing of the same properties.
	Mutex *sync.RWMutex
	// Exclude is an optional list of `sfomuseum/go-flags/multi.MultiRegexp` instances used to filter (exclude) certain properties.
	Exclude multi.MultiRegexp
	// Logger is a `log.Logger` instance used to log state and feedback.
	Logger *log.Logger
}

// EmitterCallbackFunc() returns a custom `whosonfirst/go-whosonfirst-iterate/v2/emitter.EmitterCallbackFunc` callback function
// to be invoked when iterating through Who's On First data sources that will ensure there is a corresponding "properties" JSON
// file for each of the properties in every document that is encountered.
func EmitterCallbackFunc(opts *EmitterCallbackFuncOptions) emitter.EmitterCallbackFunc {

	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Unable to load %s, because %s", path, err)
		}

		pr := gjson.GetBytes(body, "properties")

		if !pr.Exists() {
			return fmt.Errorf("%s is missing a properties dictionary!", path)
		}

		// PLEASE FOR TO go func() ME...

		for k, _ := range pr.Map() {

			_, exists := opts.Lookup.Load(k)

			if exists {
				continue
			}

			p, err := properties.NewPropertyFromKey(k)

			if err != nil {
				opts.Logger.Printf("failed to parse key (%s) for %s\n", k, path)
				continue
			}

			if p.IsName() {

				if opts.Debug {
					opts.Logger.Printf("%s is a name property, skipping\n", p)
				}

				continue
			}

			if len(opts.Exclude) > 0 {

				include := true

				for _, re := range opts.Exclude {

					if re.MatchString(p.String()) {
						include = false
						break
					}
				}

				if !include {
					continue
				}
			}

			// START OF should be updated to use gocloud.dev/blob
			// and sfomuseum/go-atomicwrite - or possibly whosonfirst/go-writer
			// which would allow writes directly back to github but it's also
			// not clear that's really necessary or desired...

			rel_path := p.RelPath()
			abs_path := filepath.Join(opts.Root, rel_path)

			opts.Mutex.Lock()

			_, err = os.Stat(abs_path)

			if os.IsNotExist(err) {

				if opts.Debug {
					opts.Logger.Printf("create %s but debugging is enabled, so don't\n", abs_path)
				} else {

					err = p.EnsureId()

					if err != nil {
						return fmt.Errorf("failed to ensure ID for %s, because %v", abs_path, err)
					}

					err = p.Write(opts.Root)

					if err != nil {
						opts.Logger.Printf("failed to write (%s) for %s, because %v\n", abs_path, path, err)
						return nil
					}
				}
			}

			// END OF should be updated to use gocloud.dev/blob

			opts.Lookup.Store(k, true)
			opts.Mutex.Unlock()
		}

		return nil
	}

	return cb
}
