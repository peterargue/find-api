package command_test

import (
	"testing"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
)

type stubResult struct{}

func (s *stubResult) String() string   { return "hello text" }
func (s *stubResult) Oneliner() string { return "hello inline" }
func (s *stubResult) JSON() any        { return map[string]any{"msg": "hello"} }

func TestFormatResultText(t *testing.T) {
	out, err := command.FormatResult(&stubResult{}, "", "text")
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello text" {
		t.Errorf("got %q, want %q", out, "hello text")
	}
}

func TestFormatResultJSON(t *testing.T) {
	out, err := command.FormatResult(&stubResult{}, "", "json")
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"msg":"hello"}` {
		t.Errorf("got %q", out)
	}
}

func TestFormatResultInline(t *testing.T) {
	out, err := command.FormatResult(&stubResult{}, "", "inline")
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello inline" {
		t.Errorf("got %q, want %q", out, "hello inline")
	}
}

func TestFormatResultFilter(t *testing.T) {
	out, err := command.FormatResult(&stubResult{}, "msg", "text")
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello" {
		t.Errorf("got %q, want %q", out, "hello")
	}
}
