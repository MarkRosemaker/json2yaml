package json2yaml_test

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"testing"

	"github.com/MarkRosemaker/json2yaml"
	"github.com/go-json-experiment/json/jsontext"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed example.json
	exampleJSON jsontext.Value
	//go:embed example.yaml
	exampleYAML []byte
)

func TestFromJSON(t *testing.T) {
	t.Parallel()

	n, err := json2yaml.Convert(exampleJSON)
	if err != nil {
		t.Fatal(err)
	}

	got, err := yaml.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, exampleYAML) {
		t.Fatalf("got: %q, want: %q", got, exampleYAML)
	}
}

func TestFromJSON_Error(t *testing.T) {
	t.Parallel()

	t.Run("empty JSON", func(t *testing.T) {
		if _, err := json2yaml.Convert(jsontext.Value(``)); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, io.EOF) {
			t.Fatalf("got: %q, want: %q", err.Error(), io.EOF)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		if _, err := json2yaml.Convert(jsontext.Value(`[{`)); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, io.EOF) {
			t.Fatalf("got: %q, want: %q", err.Error(), io.EOF)
		}

		if _, err := json2yaml.Convert(jsontext.Value(`[{}] {`)); err == nil {
			t.Fatal("expected error")
		} else if want := `expected EOF, got {`; want != err.Error() {
			t.Fatalf("got: %q, want: %q", err, io.EOF)
		}
	})

	t.Run("invalid mapping key", func(t *testing.T) {
		if _, err := json2yaml.Convert(jsontext.Value(`[{{}}]`)); err == nil {
			t.Fatal("expected error")
		} else if want := `unexpected kind for mapping key: {`; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})

	t.Run("invalid JSON map", func(t *testing.T) {
		synErr := &jsontext.SyntacticError{}
		if _, err := json2yaml.Convert(jsontext.Value(`{"foo":`)); err == nil {
			t.Fatal("expected error")
		} else if !errors.As(err, &synErr) {
			t.Fatalf("got: %T, want: %T", err, synErr)
		} else if synErr.JSONPointer != "/foo" {
			t.Fatalf("got: %q", synErr.JSONPointer)
		} else if want := `unexpected EOF`; synErr.Err.Error() != want {
			t.Fatalf("got: %q, want: %q", err, want)
		}
	})
}
