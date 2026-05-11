package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lucap123/scope/pkg/session"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <name>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a session",
	Example: `  scope delete old-target
  scope delete old-target --yes`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		yes, _ := cmd.Flags().GetBool("yes")

		// Verify it exists before asking for confirmation.
		if _, err := session.LoadSession(name); err != nil {
			return err
		}

		if !yes {
			fmt.Printf("Delete session '%s'? [y/N] ", name)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "y" && input != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		if err := session.DeleteSession(name); err != nil {
			return fmt.Errorf("failed to delete session: %w", err)
		}

		// If the deleted session was active, clear the active file.
		active, _ := session.GetActiveSessionName()
		if active == name {
			_ = session.SetActiveSession("")
			fmt.Printf("Deleted session '%s' (was active — no active session now)\n", name)
		} else {
			fmt.Printf("Deleted session '%s'\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
}
