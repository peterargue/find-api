package command_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
)

func TestSaveAndLoadToken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")

	exp := time.Now().Add(time.Hour).Unix()
	err := command.SaveToken(path, "tok123", exp)
	if err != nil {
		t.Fatalf("SaveToken: %v", err)
	}

	tok, expLoaded, err := command.LoadToken(path)
	if err != nil {
		t.Fatalf("LoadToken: %v", err)
	}
	if tok != "tok123" {
		t.Errorf("token: got %q, want %q", tok, "tok123")
	}
	if expLoaded != exp {
		t.Errorf("exp: got %d, want %d", expLoaded, exp)
	}
}

func TestLoadTokenMissing(t *testing.T) {
	_, _, err := command.LoadToken("/nonexistent/token.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadTokenExpired(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")
	exp := time.Now().Add(30 * time.Second).Unix() // expires before the 1-minute buffer
	if err := command.SaveToken(path, "old", exp); err != nil {
		t.Fatal(err)
	}
	_, _, err := command.LoadToken(path)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestTokenPath(t *testing.T) {
	p := command.TokenPath()
	if p == "" {
		t.Error("TokenPath should not be empty")
	}
	if filepath.Base(p) != "token.json" {
		t.Errorf("unexpected filename in path: %s", p)
	}
}

func TestDefaultTokenPathUsesHomeConfig(t *testing.T) {
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "find-cli", "token.json")
	got := command.TokenPath()
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
