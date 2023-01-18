package properties

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/facebookgo/atomicfile"
	"github.com/tidwall/pretty"
	"github.com/whosonfirst/go-whosonfirst-id"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// re_name is a regular expression for testing whether a property key is a `name:*` property.
var re_name *regexp.Regexp

func init() {
	re_name = regexp.MustCompile(`.*_x_.*`)
}

type PropertyType interface{}

// type Property is a struct that maps to a machine-readable data file describing a Who's On First property.
type Property struct {
	// The unique ID of this property
	Id int64 `json:"id"`
	// The name of the property
	Name string `json:"name"`
	// The namespace (prefix) of this property
	Prefix string `json:"prefix"`
	// A description of the property targeted at humans (rather than machines)
	Description string `json:"description"`
	// The expected (JSON schema) type of this property
	Type PropertyType `json:"type"`
}

// String returns the fully-qualified name (prefix + ":" + name) of this property
func (p *Property) String() string {
	return fmt.Sprintf("%s:%s", p.Prefix, p.Name)
}

// Filename() returns the filename for the serialized representation of this property.
func (p *Property) Filename() string {
	return fmt.Sprintf("%s.json", p.Name)
}

// Filename() returns the relative path (inclusive of filename) for the serialized representation of this property.
func (p *Property) RelPath() string {
	return filepath.Join(p.Prefix, p.Filename())
}

// IsName() returns a boolean value indicating whether or not the property is a `name:*` property
func (p *Property) IsName() bool {
	return re_name.MatchString(p.Name)
}

// EnsureId() ensures that the property has a unique (64-bit) identifier.
func (p *Property) EnsureId() error {

	if p.Id != -1 {
		return nil
	}

	ctx := context.Background()
	id_provider, err := id.NewProvider(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create new ID provider, %w", err)
	}

	i, err := id_provider.NewID(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create new ID, %w", err)
	}

	p.Id = i
	return nil
}

// PLEASE GIVE ME A BETTER NAME...
// (20180222/thisisaaronland)

// Write() serializes this property and writes it to 'dest' which is expected to be a directory.
func (p *Property) Write(dest string) error {

	abs_path := filepath.Join(dest, p.RelPath())
	root := filepath.Dir(abs_path)

	_, err := os.Stat(root)

	if os.IsNotExist(err) {

		err = os.MkdirAll(root, 0755)

		if err != nil {
			return fmt.Errorf("Failed to create root directory (%s), %w", root, err)
		}
	}

	enc, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("Failed to serialize property, %w", err)
	}

	fh, err := atomicfile.New(abs_path, 0644)

	if err != nil {
		return fmt.Errorf("Failed to create temporary file for writing, %w", err)
	}

	_, err = fh.Write(pretty.Pretty(enc))

	if err != nil {
		fh.Abort()
		return fmt.Errorf("Failed to write property, %w", err)
	}

	return fh.Close()
}

// NewPropertyFromKey() parses 'k' and returns a new `Property` instance.
func NewPropertyFromKey(k string) (*Property, error) {

	// PLEASE ACCOUNT FOR THINGS LIKE "src:lbl:centroid"
	// THAT OR PURGE THOSE KEYS FROM THE DATA...
	// (20180222/thisisaaronland)

	parts := strings.Split(k, ":")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid key")
	}

	p := Property{
		Id:          -1,
		Name:        parts[1],
		Prefix:      parts[0],
		Description: "",
		Type:        "",
	}

	return &p, nil
}

// NewPropertyFromFiles() reads and parses the contents of 'path' and returns a new `Property` instance.
func NewPropertyFromFile(path string) (*Property, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, fmt.Errorf("Failed to determine absolute path for file, %w", err)
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to open '%s' for reading, %w", abs_path, err)
	}

	defer fh.Close()

	return NewPropertyFromReader(fh)
}

// NewPropertyFromFiles() reads and parses the contents of 'fh' and returns a new `Property` instance.
func NewPropertyFromReader(fh io.Reader) (*Property, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, fmt.Errorf("Failed to read contents of reader, %w", err)
	}

	var prop Property

	err = json.Unmarshal(body, &prop)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode property, %w", err)
	}

	return &prop, nil
}
