package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lucapnss/scope/pkg/session"
)

// RunTool executes an external tool with session headers/proxy injected.
func RunTool(toolName string, args []string, verbose bool) error {
	s, err := session.GetActiveSession()
	if err != nil {
		return err
	}

	var cmdArgs []string

	switch toolName {
	case "curl":
		cmdArgs = buildCurlArgs(s, args)
	case "ffuf":
		cmdArgs = buildFfufArgs(s, args)
	case "http":
		cmdArgs = buildHttpieArgs(s, args)
	case "nuclei":
		cmdArgs = buildNucleiArgs(s, args)
	case "sqlmap":
		cmdArgs = buildSqlmapArgs(s, args)
	default:
		return fmt.Errorf("tool '%s' is not supported", toolName)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "[scope] %s %s\n", toolName, strings.Join(cmdArgs, " "))
	}

	cmd := exec.Command(toolName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// Don't wrap exec.ExitError — let the exit code propagate naturally.
		if _, ok := err.(*exec.ExitError); ok {
			os.Exit(cmd.ProcessState.ExitCode())
		}
		return fmt.Errorf("failed to run %s: %w", toolName, err)
	}
	return nil
}

// BuildArgs returns the argument slice for a tool without running it.
// Used by 'scope export <tool>'.
func BuildArgs(toolName string, s *session.Session, args []string) ([]string, error) {
	switch toolName {
	case "curl":
		return buildCurlArgs(s, args), nil
	case "ffuf":
		return buildFfufArgs(s, args), nil
	case "http":
		return buildHttpieArgs(s, args), nil
	case "nuclei":
		return buildNucleiArgs(s, args), nil
	case "sqlmap":
		return buildSqlmapArgs(s, args), nil
	default:
		return nil, fmt.Errorf("tool '%s' is not supported", toolName)
	}
}

// SupportedTools returns the list of tools scope can wrap.
func SupportedTools() []string {
	return []string{"curl", "ffuf", "http", "nuclei", "sqlmap"}
}

// resolveURL prepends the session base URL when the arg starts with '/'.
func resolveURL(base, arg string) string {
	if strings.HasPrefix(arg, "/") {
		return strings.TrimRight(base, "/") + arg
	}
	return arg
}

func buildCurlArgs(s *session.Session, args []string) []string {
	var out []string

	for k, v := range s.Headers {
		out = append(out, "-H", fmt.Sprintf("%s: %s", k, v))
	}
	if s.Proxy != "" {
		out = append(out, "--proxy", s.Proxy)
	}

	for _, arg := range args {
		out = append(out, resolveURL(s.URL, arg))
	}
	return out
}

func buildFfufArgs(s *session.Session, args []string) []string {
	var out []string

	for k, v := range s.Headers {
		out = append(out, "-H", fmt.Sprintf("%s: %s", k, v))
	}
	if s.Proxy != "" {
		out = append(out, "-x", s.Proxy)
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "-u" && i+1 < len(args) {
			out = append(out, "-u", resolveURL(s.URL, args[i+1]))
			i++
		} else {
			out = append(out, args[i])
		}
	}
	return out
}

func buildHttpieArgs(s *session.Session, args []string) []string {
	var out []string

	// httpie proxy flag: --proxy http:http://host:port
	if s.Proxy != "" {
		out = append(out, "--proxy", "http:"+s.Proxy)
	}

	for _, arg := range args {
		out = append(out, resolveURL(s.URL, arg))
	}

	// httpie headers come after the URL as Header:Value pairs
	for k, v := range s.Headers {
		out = append(out, fmt.Sprintf("%s:%s", k, v))
	}
	return out
}

func buildNucleiArgs(s *session.Session, args []string) []string {
	var out []string

	for k, v := range s.Headers {
		out = append(out, "-H", fmt.Sprintf("%s: %s", k, v))
	}
	if s.Proxy != "" {
		out = append(out, "-proxy", s.Proxy)
	}

	// If no -u flag was provided, inject the session URL automatically.
	hasTarget := false
	for i := 0; i < len(args); i++ {
		if args[i] == "-u" && i+1 < len(args) {
			out = append(out, "-u", resolveURL(s.URL, args[i+1]))
			i++
			hasTarget = true
		} else {
			out = append(out, args[i])
		}
	}
	if !hasTarget && s.URL != "" {
		out = append(out, "-u", s.URL)
	}

	return out
}

func buildSqlmapArgs(s *session.Session, args []string) []string {
	var out []string

	// sqlmap accepts multiple headers via --headers with \n separator
	var hdrs []string
	for k, v := range s.Headers {
		hdrs = append(hdrs, fmt.Sprintf("%s: %s", k, v))
	}
	if len(hdrs) > 0 {
		out = append(out, "--headers", strings.Join(hdrs, "\n"))
	}
	if s.Proxy != "" {
		out = append(out, "--proxy", s.Proxy)
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "-u" && i+1 < len(args) {
			out = append(out, "-u", resolveURL(s.URL, args[i+1]))
			i++
		} else {
			out = append(out, args[i])
		}
	}
	return out
}
