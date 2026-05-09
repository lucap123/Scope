package cmd

import (
	"fmt"
	"sort"

	"github.com/lucapnss/scope/pkg/session"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show session details",
	Example: `  scope show               # show active session
  scope show target        # show a specific session`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reveal, _ := cmd.Flags().GetBool("reveal")

		var name string
		var err error
		if len(args) == 1 {
			name = args[0]
		} else {
			name, err = session.GetActiveSessionName()
			if err != nil || name == "" {
				return fmt.Errorf("no active session — specify a name or run 'scope use <name>'")
			}
		}

		s, err := session.LoadSession(name)
		if err != nil {
			return err
		}

		activeName, _ := session.GetActiveSessionName()
		activeTag := ""
		if s.Name == activeName {
			activeTag = "  (active)"
		}

		fmt.Printf("\n  Session:   %s%s\n", s.Name, activeTag)
		fmt.Printf("  URL:       %s\n", s.URL)

		// Headers — sorted for deterministic output, secrets masked by default.
		fmt.Println("  Headers:")
		if len(s.Headers) == 0 {
			fmt.Println("             —")
		} else {
			keys := make([]string, 0, len(s.Headers))
			for k := range s.Headers {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := s.Headers[k]
				if !reveal {
					v = session.MaskValue(v)
				}
				fmt.Printf("             %s: %s\n", k, v)
			}
		}

		if s.Proxy != "" {
			fmt.Printf("  Proxy:     %s\n", s.Proxy)
		} else {
			fmt.Println("  Proxy:     —")
		}

		fmt.Printf("  Created:   %s\n", s.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Printf("  Updated:   %s\n", s.UpdatedAt.Format("2006-01-02 15:04"))

		if len(s.Notes) > 0 {
			fmt.Println("  Notes:")
			for _, n := range s.Notes {
				fmt.Printf("    [%s] %s\n", n.Timestamp.Format("2006-01-02 15:04"), n.Content)
			}
		} else {
			fmt.Println("  Notes:     —")
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().Bool("reveal", false, "Show full header values (not masked)")
}
