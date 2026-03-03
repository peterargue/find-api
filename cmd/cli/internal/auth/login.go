package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	findapi "github.com/peterargue/find-api"
	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCobra = &cobra.Command{
	Use:   "login",
	Short: "Log in and store a 7-day authentication token",
	Long: `Prompts for your FindLabs username and password, generates a 7-day
token, and stores it at ~/.config/find-cli/token.json.`,
}

func runLogin(args []string, flags *command.GlobalFlags) (command.Result, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading username: %w", err)
	}
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return nil, fmt.Errorf("reading password: %w", err)
	}
	password := strings.TrimSpace(string(passwordBytes))

	client := findapi.NewClient(username, password)
	resp, err := client.Auth.GenerateToken(context.Background(), 168*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	if err := command.SaveToken(command.TokenPath(), resp.AccessToken, resp.Exp); err != nil {
		return nil, fmt.Errorf("saving token: %w", err)
	}

	expiry := time.Unix(resp.Exp, 0).Format(time.RFC3339)
	return &loginResult{expiry: expiry}, nil
}
