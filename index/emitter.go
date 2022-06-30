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

type EmitterCallbackFuncOptions struct {
	Debug   bool
	Lookup  *sync.Map
	Root    string
	Mutex   *sync.RWMutex
	Exclude multi.MultiRegexp
	Logger  *log.Logger
}

func EmitterCallbackFunc(opts *EmitterCallbackFuncOptions) emitter.EmitterCallbackFunc {

	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		opts.Logger.Println(path)

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
			defer opts.Mutex.Unlock()

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
		}

		return nil

	}

	return cb
}
