package cmd

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lucap123/scope/pkg/session"
	"github.com/spf13/cobra"
)

// BurpItems is the top-level element in a Burp Suite XML export.
type BurpItems struct {
	XMLName xml.Name   `xml:"items"`
	Items   []BurpItem `xml:"item"`
}

// BurpItem represents a single request/response pair in the Burp export.
type BurpItem struct {
	URL     string `xml:"url"`
	Request struct {
		Base64  bool   `xml:"base64,attr"`
		Content string `xml:",chardata"`
	} `xml:"request"`
}

// skipHeaders are headers we don't want to import — they change per-request.
var skipHeaders = map[string]bool{
	"host":           true,
	"content-length": true,
	"connection":     true,
	"accept-encoding": true,
}

var importCmd = &cobra.Command{
	Use:   "import <type> <file>",
	Short: "Import a session from a file",
	Long: `Import a session from an external tool's export file.

Supported types:
  burp   — Burp Suite XML export (right-click request > Save item)`,
	Example: `  scope import burp burp-export.xml`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		importType := args[0]
		filePath := args[1]

		switch strings.ToLower(importType) {
		case "burp":
			return importBurp(filePath)
		default:
			return fmt.Errorf("unknown import type '%s' — supported types: burp", importType)
		}
	},
}

func importBurp(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	var items BurpItems
	if err := xml.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("invalid Burp XML: %w", err)
	}

	if len(items.Items) == 0 {
		return fmt.Errorf("no request items found in Burp export")
	}

	item := items.Items[0]

	// Decode the HTTP request (may be base64 encoded).
	var rawRequest []byte
	if item.Request.Base64 {
		rawRequest, err = base64.StdEncoding.DecodeString(strings.TrimSpace(item.Request.Content))
		if err != nil {
			return fmt.Errorf("failed to decode request data: %w", err)
		}
	} else {
		rawRequest = []byte(item.Request.Content)
	}

	// Derive base URL from the item URL.
	baseURL := item.URL
	if parsed, err := url.Parse(item.URL); err == nil {
		baseURL = fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
	}

	s := &session.Session{
		Name:      "burp-" + time.Now().Format("20060102-150405"),
		URL:       baseURL,
		Headers:   make(map[string]string),
		CreatedAt: time.Now(),
	}

	// Parse HTTP request headers (skip the request line and stop at blank line).
	lines := strings.Split(strings.ReplaceAll(string(rawRequest), "\r\n", "\n"), "\n")
	for i, line := range lines {
		if i == 0 {
			continue // skip "GET /path HTTP/1.1"
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if skipHeaders[strings.ToLower(key)] {
			continue
		}
		s.Headers[key] = val
	}

	if err := session.SaveSession(s); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	fmt.Printf("Imported session '%s'\n", s.Name)
	fmt.Printf("  URL:     %s\n", s.URL)
	fmt.Printf("  Headers: %d imported\n", len(s.Headers))
	fmt.Printf("  Rename with: scope set %s url <url>\n", s.Name)
	return nil
}

func init() {
	rootCmd.AddCommand(importCmd)
}
