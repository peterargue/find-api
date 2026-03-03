package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	findapi "github.com/peterargue/find-api"
)

type tokenFile struct {
	AccessToken string `json:"access_token"`
	Exp         int64  `json:"exp"`
}

// TokenPath returns the path to the stored token file.
func TokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "find-cli", "token.json")
}

// SaveToken writes a token and expiry to the given path, creating parent dirs.
func SaveToken(path, token string, exp int64) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create token dir: %w", err)
	}
	b, err := json.Marshal(tokenFile{AccessToken: token, Exp: exp})
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

// LoadToken reads and validates a token from the given path.
// Returns an error if the file is missing or the token is expired.
func LoadToken(path string) (token string, exp int64, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", 0, fmt.Errorf("no stored token: %w", err)
	}
	var tf tokenFile
	if err := json.Unmarshal(b, &tf); err != nil {
		return "", 0, fmt.Errorf("invalid token file: %w", err)
	}
	if time.Now().Add(time.Minute).After(time.Unix(tf.Exp, 0)) {
		return "", 0, errors.New("stored token is expired or expiring soon")
	}
	return tf.AccessToken, tf.Exp, nil
}

// MustLoadClient loads the stored token and returns a configured API client.
// If the token is missing or expired it prints a helpful message and exits.
func MustLoadClient() *findapi.Client {
	token, exp, err := LoadToken(TokenPath())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Not authenticated. Run: find auth login")
		os.Exit(1)
	}
	return findapi.NewClient("", "", findapi.WithToken(token, exp))
}
