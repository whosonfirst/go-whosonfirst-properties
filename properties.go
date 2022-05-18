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

var re_name *regexp.Regexp

func init() {
	re_name = regexp.MustCompile(`.*_x_.*`)
}

type Property struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func (p *Property) String() string {
	return fmt.Sprintf("%s:%s", p.Prefix, p.Name)
}

func (p *Property) Filename() string {
	return fmt.Sprintf("%s.json", p.Name)
}

func (p *Property) RelPath() string {
	return filepath.Join(p.Prefix, p.Filename())
}

func (p *Property) IsName() bool {
	return re_name.MatchString(p.Name)
}

func (p *Property) EnsureId() error {

	if p.Id != -1 {
		return nil
	}

	ctx := context.Background()
	id_provider, err := id.NewProvider(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create new ID provider, %w", err)
	}

	i, err := id_provider.NewID()

	if err != nil {
		return fmt.Errorf("Failed to create new ID, %w", err)
	}

	p.Id = i
	return nil
}

// PLEASE GIVE ME A BETTER NAME...
// (20180222/thisisaaronland)

func (p *Property) Write(dest string) error {

	abs_path := filepath.Join(dest, p.RelPath())
	root := filepath.Dir(abs_path)

	_, err := os.Stat(root)

	if os.IsNotExist(err) {

		err = os.MkdirAll(root, 0755)

		if err != nil {
			return err
		}
	}

	enc, err := json.Marshal(p)

	if err != nil {
		return err
	}

	fh, err := atomicfile.New(abs_path, 0644)

	if err != nil {
		return err
	}

	_, err = fh.Write(pretty.Pretty(enc))

	if err != nil {
		fh.Abort()
		return err
	}

	return fh.Close()
}

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

func NewPropertyFromFile(path string) (*Property, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return NewPropertyFromReader(fh)
}

func NewPropertyFromReader(fh io.Reader) (*Property, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	var prop Property

	err = json.Unmarshal(body, &prop)

	if err != nil {
		return nil, err
	}

	return &prop, nil
}
