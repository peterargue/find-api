package auth

import (
	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/spf13/cobra"
)

// Cmd is the "auth" subcommand group.
var Cmd = &cobra.Command{
	Use:              "auth",
	Short:            "Authentication commands",
	TraverseChildren: true,
}

func init() {
	loginCmd.AddToParent(Cmd)
}

// loginResult is the Result for a successful login.
type loginResult struct {
	expiry string
}

func (r *loginResult) String() string   { return "Logged in. Token valid until: " + r.expiry }
func (r *loginResult) Oneliner() string { return r.expiry }
func (r *loginResult) JSON() any        { return map[string]any{"token_expiry": r.expiry} }

var loginCmd = &command.Command{
	Cmd:   loginCobra,
	Flags: loginFlagsVal,
	Run:   runLogin,
}
