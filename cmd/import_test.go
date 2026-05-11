package cmd

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/lucap123/scope/pkg/session"
)

func TestImportBurp(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "scope-import-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)
	t.Setenv("SCOPE_HOME", tempHome)

	request := "GET /api/v1/users HTTP/1.1\r\nHost: example.com\r\nAuthorization: Bearer test-token\r\nContent-Length: 0\r\n\r\n"
	requestEncoded := base64.StdEncoding.EncodeToString([]byte(request))

	xmlContent := `<?xml version="1.0"?>
<items>
  <item>
    <url>https://example.com/api/v1/users</url>
    <request base64="true">` + requestEncoded + `</request>
  </item>
</items>`

	xmlFile := tempHome + "/test-burp.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlContent), 0600); err != nil {
		t.Fatal(err)
	}

	if err := importBurp(xmlFile); err != nil {
		t.Fatalf("importBurp failed: %v", err)
	}

	sessions, err := session.ListSessions()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	s, err := session.LoadSession(sessions[0])
	if err != nil {
		t.Fatal(err)
	}

	// importBurp now stores only the base URL (scheme + host).
	want := "https://example.com"
	if s.URL != want {
		t.Errorf("URL: want %s, got %s", want, s.URL)
	}

	if s.Headers["Authorization"] != "Bearer test-token" {
		t.Errorf("Authorization header: want 'Bearer test-token', got %q", s.Headers["Authorization"])
	}

	// Host and Content-Length must be skipped.
	if _, ok := s.Headers["Host"]; ok {
		t.Error("Host header should not be imported")
	}
	if _, ok := s.Headers["Content-Length"]; ok {
		t.Error("Content-Length header should not be imported")
	}
}

func TestImportBurpEmptyFile(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "scope-import-test-*")
	defer os.RemoveAll(tempHome)
	t.Setenv("SCOPE_HOME", tempHome)

	xmlContent := `<?xml version="1.0"?><items></items>`
	xmlFile := tempHome + "/empty.xml"
	os.WriteFile(xmlFile, []byte(xmlContent), 0600)

	err := importBurp(xmlFile)
	if err == nil {
		t.Error("expected error for empty Burp file, got nil")
	}
}
