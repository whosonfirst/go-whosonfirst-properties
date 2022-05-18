package properties

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPropertiesFromKey(t *testing.T) {

	pr, err := NewPropertyFromKey("hello:world")

	if err != nil {
		t.Fatalf("Failed to create new property from key, %v", err)
	}

	if pr.Name != "world" {
		t.Fatalf("Invalid name: %s", pr.Name)
	}

	if pr.Prefix != "hello" {
		t.Fatalf("Invalid prefix: %s", pr.Prefix)
	}

	if pr.IsName() {
		t.Fatalf("Property reported as 'name' property")
	}

	if pr.RelPath() != "hello/world.json" {
		t.Fatalf("Invalid rel path: %s", pr.RelPath())
	}
}

func TestNamePropertiesFromKey(t *testing.T) {

	pr, err := NewPropertyFromKey("name:eng_x_example")

	if err != nil {
		t.Fatalf("Failed to create new property from key, %v", err)
	}

	if pr.Name != "eng_x_example" {
		t.Fatalf("Invalid name: %s", pr.Name)
	}

	if pr.Prefix != "name" {
		t.Fatalf("Invalid prefix: %s", pr.Prefix)
	}

	if !pr.IsName() {
		t.Fatalf("Expected property to be name")
	}
}

func TestPropertiesFromFile(t *testing.T) {

	rel_path := "fixtures/wof/abbreviation.json"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for '%s', %v", rel_path, err)
	}

	pr, err := NewPropertyFromFile(abs_path)

	if err != nil {
		t.Fatalf("Failed to create property for '%s', %v", abs_path, err)
	}

	if pr.Id != 1158807929 {
		t.Fatalf("Invalid ID: %d", pr.Id)
	}

	if pr.Name != "abbreviation" {
		t.Fatalf("Invalid name: %s", pr.Name)
	}

	if pr.Prefix != "wof" {
		t.Fatalf("Invalid prefix: %s", pr.Prefix)
	}
}

func TestWriteProperties(t *testing.T) {

	pr, err := NewPropertyFromKey("hello:world")

	if err != nil {
		t.Fatalf("Failed to create new property from key, %v", err)
	}

	if pr.Name != "world" {
		t.Fatalf("Invalid name: %s", pr.Name)
	}

	if pr.Prefix != "hello" {
		t.Fatalf("Invalid prefix: %s", pr.Prefix)
	}

	if pr.IsName() {
		t.Fatalf("Property reported as 'name' property")
	}

	dir, err := ioutil.TempDir("", "properties")

	if err != nil {
		t.Fatalf("Failed to create temp dir, %v", err)
	}

	defer os.RemoveAll(dir)

	err = pr.Write(dir)

	if err != nil {
		t.Fatalf("Failed to write property, %v", err)
	}

}
