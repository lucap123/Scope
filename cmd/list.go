package cmd

import (
	"fmt"

	"github.com/lucapnss/scope/pkg/session"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessions, err := session.ListSessions()
		if err != nil {
			return fmt.Errorf("failed to list sessions: %w", err)
		}

		if len(sessions) == 0 {
			fmt.Println("No sessions yet. Create one with: scope new <name>")
			return nil
		}

		active, _ := session.GetActiveSessionName()

		fmt.Println()
		for _, name := range sessions {
			tag := "  "
			if name == active {
				tag = "* "
			}

			s, err := session.LoadSession(name)
			url := ""
			headerCount := 0
			if err == nil {
				url = s.URL
				headerCount = len(s.Headers)
			}

			hInfo := ""
			if headerCount > 0 {
				hInfo = fmt.Sprintf("  [%d header", headerCount)
				if headerCount > 1 {
					hInfo += "s"
				}
				hInfo += "]"
			}

			fmt.Printf("  %s%-24s %-40s%s\n", tag, name, url, hInfo)
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
