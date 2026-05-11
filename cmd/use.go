package cmd

import (
	"fmt"

	"github.com/lucap123/scope/pkg/session"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:     "use <name>",
	Short:   "Set the active session",
	Example: `  scope use bugcrowd-target`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if _, err := session.LoadSession(name); err != nil {
			return fmt.Errorf("session '%s' does not exist — create it with 'scope new %s'", name, name)
		}

		if err := session.SetActiveSession(name); err != nil {
			return fmt.Errorf("failed to set active session: %w", err)
		}

		fmt.Printf("Active session: %s\n", name)
		fmt.Println("  Tip: set SCOPE_SESSION=<name> in your shell to override per-terminal.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
