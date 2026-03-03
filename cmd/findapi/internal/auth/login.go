package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	findapi "github.com/peterargue/find-api"
	"github.com/peterargue/find-api/cmd/findapi/internal/command"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCobra = &cobra.Command{
	Use:   "login",
	Short: "Log in and store a 7-day authentication token",
	Long: `Logs in to the FindLabs API and stores a 7-day token at ~/.config/find-cli/token.json.

Credentials can be provided via flags or entered interactively when omitted.`,
	Example: "find auth login\nfind auth login --username alice --password s3cr3t",
}

type loginFlags struct {
	Username string `flag:"username" info:"FindLabs username (prompted if not set)"`
	Password string `flag:"password" info:"FindLabs password (prompted if not set)"`
}

var loginFlagsVal = &loginFlags{}

// maxTokenDuration is the maximum session length the FindLabs API supports.
const maxTokenDuration = 7 * 24 * time.Hour

func runLogin(args []string, flags *command.GlobalFlags) (command.Result, error) {
	username := loginFlagsVal.Username
	password := loginFlagsVal.Password

	if username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Username: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("reading username: %w", err)
		}
		username = strings.TrimSpace(line)
	}

	if password == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return nil, fmt.Errorf("could not read password securely (is stdin a terminal?): %w", err)
		}
		password = strings.TrimSpace(string(passwordBytes))
	}

	client := findapi.NewClient(username, password)
	resp, err := client.Auth.GenerateToken(context.Background(), maxTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	if err := command.SaveToken(command.TokenPath(), resp.AccessToken, resp.Exp); err != nil {
		return nil, fmt.Errorf("saving token: %w", err)
	}

	expiry := time.Unix(resp.Exp, 0).Format(time.RFC3339)
	return &loginResult{expiry: expiry}, nil
}
