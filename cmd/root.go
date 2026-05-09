package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "scope",
	Short: "scope — CLI session manager for web hacking",
	Long: `scope keeps a named session per target so you never copy-paste
auth tokens again. Set a base URL, headers, and proxy once — then
run curl, ffuf, nuclei, sqlmap and more with everything auto-injected.

  scope new target --url https://api.example.com
  scope set target header "Authorization: Bearer <token>"
  scope use target
  scope curl /users
  scope ffuf -w wordlist.txt -u /FUZZ

Sessions are stored in ~/.scope/sessions/ as plain JSON files.
Set SCOPE_SESSION=<name> to override the active session per terminal.`,
	Version: Version,
}

// Verbose flag is available to all subcommands.
var Verbose bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print the full command before executing")
}
