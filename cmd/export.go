package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/lucapnss/scope/pkg/runner"
	"github.com/lucapnss/scope/pkg/session"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [session|tool] [args...]",
	Short: "Export a session or generate a ready-to-paste command",
	Long: `Two modes:

  1. Export a session as JSON or shell env vars:
       scope export bugcrowd-target --format json
       scope export bugcrowd-target --format env

  2. Generate a copy-paste command for a specific tool:
       scope export curl /users
       scope export curl -X POST /login -d '{"user":"test"}'
       scope export ffuf -w words.txt -u /FUZZ`,
	Example: `  scope export bugcrowd-target --format json > session.json
  scope export bugcrowd-target --format env
  scope export curl /admin/settings`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		target := args[0]

		// Check if the first arg is a supported tool name.
		for _, t := range runner.SupportedTools() {
			if t == target {
				return exportToolCommand(target, args[1:])
			}
		}

		// Otherwise treat it as a session name.
		return exportSession(target, format)
	},
}

func exportSession(name, format string) error {
	s, err := session.LoadSession(name)
	if err != nil {
		return err
	}

	switch format {
	case "json":
		data, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))

	case "env":
		fmt.Printf("export SCOPE_URL=\"%s\"\n", s.URL)
		if s.Proxy != "" {
			fmt.Printf("export SCOPE_PROXY=\"%s\"\n", s.Proxy)
		}
		// Sort headers for deterministic output.
		keys := make([]string, 0, len(s.Headers))
		for k := range s.Headers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			envKey := strings.ToUpper(strings.ReplaceAll(k, "-", "_"))
			fmt.Printf("export SCOPE_%s=\"%s\"\n", envKey, s.Headers[k])
		}

	default:
		return fmt.Errorf("unknown format '%s' — valid formats: json, env", format)
	}
	return nil
}

func exportToolCommand(toolName string, args []string) error {
	s, err := session.GetActiveSession()
	if err != nil {
		return err
	}

	builtArgs, err := runner.BuildArgs(toolName, s, args)
	if err != nil {
		return err
	}

	// Quote args that contain spaces so the output is copy-paste safe.
	quoted := make([]string, 0, len(builtArgs)+1)
	quoted = append(quoted, toolName)
	for _, a := range builtArgs {
		if strings.ContainsAny(a, " \t") {
			a = fmt.Sprintf("'%s'", a)
		}
		quoted = append(quoted, a)
	}

	fmt.Println(strings.Join(quoted, " "))
	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("format", "f", "json", "Export format: json or env")
}
