package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Session holds all state for a single hacking target.
type Session struct {
	Name      string            `json:"name"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Proxy     string            `json:"proxy"`
	Notes     []Note            `json:"notes"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Note is a timestamped free-text annotation attached to a session.
type Note struct {
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

// MaskValue truncates a sensitive value for display (e.g. tokens, cookies).
// Shows first 6 and last 4 characters, rest replaced with dots.
func MaskValue(v string) string {
	if len(v) <= 12 {
		return strings.Repeat("*", len(v))
	}
	return v[:6] + "..." + v[len(v)-4:]
}

// GetConfigDir returns (and creates if needed) the ~/.scope directory.
// Respects the SCOPE_HOME env var for overrides (useful in tests).
func GetConfigDir() (string, error) {
	dir := os.Getenv("SCOPE_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine home directory: %w", err)
		}
		dir = filepath.Join(home, ".scope")
	}

	for _, sub := range []string{dir, filepath.Join(dir, "sessions")} {
		if err := os.MkdirAll(sub, 0700); err != nil {
			return "", fmt.Errorf("cannot create config dir %s: %w", sub, err)
		}
	}
	return dir, nil
}

// sessionPath returns the path to a named session's JSON file.
func sessionPath(name string) (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sessions", name+".json"), nil
}

// SaveSession persists a session to disk.
func SaveSession(s *Session) error {
	s.UpdatedAt = time.Now()
	path, err := sessionPath(s.Name)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadSession reads and parses a session from disk.
func LoadSession(name string) (*Session, error) {
	path, err := sessionPath(name)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session '%s' does not exist", name)
		}
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("corrupt session file for '%s': %w", name, err)
	}
	return &s, nil
}

// DeleteSession removes a session file from disk.
func DeleteSession(name string) error {
	path, err := sessionPath(name)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("session '%s' does not exist", name)
		}
		return err
	}
	return nil
}

// ListSessions returns all session names sorted alphabetically.
func ListSessions() ([]string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(dir, "sessions"))
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return names, nil
}

// activeFilePath returns the path to the active-session marker file.
// Uses SCOPE_SESSION env var first, then falls back to a per-process
// file so multiple terminals can have independent active sessions.
func activeFilePath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "active"), nil
}

// SetActiveSession writes the active session name to disk.
func SetActiveSession(name string) error {
	path, err := activeFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(name), 0600)
}

// GetActiveSessionName returns the currently active session name.
// Priority order:
//  1. SCOPE_SESSION env var (per-terminal, set manually by user)
//  2. ~/.scope/active file (global fallback, set by 'scope use')
func GetActiveSessionName() (string, error) {
	// Env var takes priority — lets users do: export SCOPE_SESSION=target
	if name := os.Getenv("SCOPE_SESSION"); name != "" {
		return strings.TrimSpace(name), nil
	}

	path, err := activeFilePath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// GetActiveSession loads the currently active session.
func GetActiveSession() (*Session, error) {
	name, err := GetActiveSessionName()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("no active session — run 'scope use <name>' first")
	}
	return LoadSession(name)
}

// UnsetHeader removes a header from a session by key (case-insensitive).
func (s *Session) UnsetHeader(key string) bool {
	for k := range s.Headers {
		if strings.EqualFold(k, key) {
			delete(s.Headers, k)
			return true
		}
	}
	return false
}
