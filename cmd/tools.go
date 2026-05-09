package cmd

import (
	"github.com/lucapnss/scope/pkg/runner"
	"github.com/spf13/cobra"
)

func registerTool(name, short string) {
	c := &cobra.Command{
		Use:                name + " [args]",
		Short:              short,
		DisableFlagParsing: true, // pass all flags straight to the underlying tool
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.RunTool(name, args, Verbose)
		},
	}
	rootCmd.AddCommand(c)
}

func init() {
	registerTool("curl", "Run curl with session headers and proxy injected")
	registerTool("ffuf", "Run ffuf with session headers and proxy injected")
	registerTool("http", "Run httpie with session headers and proxy injected")
	registerTool("nuclei", "Run nuclei with session headers and proxy injected")
	registerTool("sqlmap", "Run sqlmap with session headers and proxy injected")
}
