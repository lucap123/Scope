package cmd

import (
	"fmt"
	"time"

	"github.com/lucap123/scope/pkg/session"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new session",
	Example: `  scope new bugcrowd-target
  scope new bugcrowd-target --url https://api.target.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		url, _ := cmd.Flags().GetString("url")

		// Don't overwrite an existing session silently.
		if _, err := session.LoadSession(name); err == nil {
			return fmt.Errorf("session '%s' already exists — use 'scope set' to update it", name)
		}

		s := &session.Session{
			Name:      name,
			URL:       url,
			Headers:   make(map[string]string),
			CreatedAt: time.Now(),
		}

		if err := session.SaveSession(s); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		fmt.Printf("Created session: %s\n", name)
		if url == "" {
			fmt.Printf("  Set a URL with: scope set %s url <url>\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("url", "u", "", "Base URL for the session")
}
