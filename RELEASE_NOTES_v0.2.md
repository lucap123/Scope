## What's new in v0.2

### New features
- **Burp Suite import** — right-click a request in Burp → Save item → `scope import burp burp-export.xml`. Automatically parses base URL, auth headers, and cookies.
- **Session export** — export any session as JSON or env format: `scope export <name> --format json`
- **Command export** — generate copy-paste ready commands without running them: `scope export curl /users`
- **Notes** — attach recon notes to your active session: `scope note "Admin panel at /manage no auth check"`

### What was already in v0.1
- **Session management** (new, use, list, delete, show)
- **Auto-inject** headers, proxy, and base URL into curl, ffuf, nuclei, sqlmap, and more.
- **Verbose mode** — use `--verbose` to preview full commands before they run.
- **Local-first storage** — plain JSON files in `~/.scope/sessions/`. No cloud, no telemetry.

## Install

```bash
go install github.com/lucap123/scope@latest
```

Or build from source:
```bash
git clone https://github.com/lucap123/scope
cd scope && make install
```
