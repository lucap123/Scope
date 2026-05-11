package cmd

import (
	"fmt"
	"time"

	"github.com/lucap123/scope/pkg/session"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note <content>",
	Short: "Add a note to the active session",
	Example: `  scope note "Login endpoint: /api/v2/auth"
  scope note "Admin panel at /manage — no auth check!"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := session.GetActiveSession()
		if err != nil {
			return err
		}

		s.Notes = append(s.Notes, session.Note{
			Timestamp: time.Now(),
			Content:   args[0],
		})

		if err := session.SaveSession(s); err != nil {
			return fmt.Errorf("failed to save note: %w", err)
		}

		fmt.Println("Note added.")
		return nil
	},
}

var notesCmd = &cobra.Command{
	Use:   "notes [session]",
	Short: "List notes for a session",
	Example: `  scope notes                  # notes for active session
  scope notes bugcrowd-target  # notes for a specific session`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var s *session.Session
		var err error

		if len(args) == 1 {
			s, err = session.LoadSession(args[0])
		} else {
			s, err = session.GetActiveSession()
		}
		if err != nil {
			return err
		}

		if len(s.Notes) == 0 {
			fmt.Printf("No notes for session '%s'.\n", s.Name)
			return nil
		}

		fmt.Printf("\n  Notes for %s:\n\n", s.Name)
		for _, n := range s.Notes {
			fmt.Printf("  [%s] %s\n", n.Timestamp.Format("2006-01-02 15:04"), n.Content)
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
	rootCmd.AddCommand(notesCmd)
}
