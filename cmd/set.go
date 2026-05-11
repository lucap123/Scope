package cmd

import (
	"fmt"
	"strings"

	"github.com/lucap123/scope/pkg/session"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [session] <key> <value>",
	Short: "Set a property on a session",
	Long: `Set a property on a session. If session name is omitted, uses the active session.

Valid keys:
  url      — base URL for the target
  header   — add or replace a header (format: "Key: Value")
  proxy    — proxy URL (e.g. http://127.0.0.1:8080)`,
	Example: `  scope set target url https://api.target.com
  scope set target header "Authorization: Bearer eyJ..."
  scope set target header "Cookie: session=abc123"
  scope set target proxy http://127.0.0.1:8080
  scope set url https://api.target.com        # uses active session`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, key, value string

		if len(args) == 3 {
			name = args[0]
			key = args[1]
			value = args[2]
		} else {
			var err error
			name, err = session.GetActiveSessionName()
			if err != nil || name == "" {
				return fmt.Errorf("no active session — specify a name or run 'scope use <name>'")
			}
			key = args[0]
			value = args[1]
		}

		s, err := session.LoadSession(name)
		if err != nil {
			return err
		}

		switch strings.ToLower(key) {
		case "url":
			s.URL = value
			fmt.Printf("Set url = %s\n", value)

		case "header":
			parts := strings.SplitN(value, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("header must be in 'Key: Value' format, got: %s", value)
			}
			hKey := strings.TrimSpace(parts[0])
			hVal := strings.TrimSpace(parts[1])
			if s.Headers == nil {
				s.Headers = make(map[string]string)
			}
			s.Headers[hKey] = hVal
			fmt.Printf("Set header %s\n", hKey)

		case "proxy":
			s.Proxy = value
			fmt.Printf("Set proxy = %s\n", value)

		default:
			return fmt.Errorf("unknown key '%s' — valid keys: url, header, proxy", key)
		}

		if err := session.SaveSession(s); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}
		return nil
	},
}

var unsetCmd = &cobra.Command{
	Use:   "unset [session] header <key>",
	Short: "Remove a header from a session",
	Long:  `Remove a specific header from a session by header name (case-insensitive).`,
	Example: `  scope unset header Authorization
  scope unset target header Cookie`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, key string

		if len(args) == 3 {
			name = args[0]
			if strings.ToLower(args[1]) != "header" {
				return fmt.Errorf("only 'header' can be unset, got: %s", args[1])
			}
			key = args[2]
		} else {
			var err error
			name, err = session.GetActiveSessionName()
			if err != nil || name == "" {
				return fmt.Errorf("no active session — specify a name or run 'scope use <name>'")
			}
			if strings.ToLower(args[0]) != "header" {
				return fmt.Errorf("only 'header' can be unset, got: %s", args[0])
			}
			key = args[1]
		}

		s, err := session.LoadSession(name)
		if err != nil {
			return err
		}

		if !s.UnsetHeader(key) {
			return fmt.Errorf("header '%s' not found in session '%s'", key, name)
		}

		if err := session.SaveSession(s); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}
		fmt.Printf("Removed header: %s\n", key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(unsetCmd)
}
