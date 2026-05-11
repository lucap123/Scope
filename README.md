# scope
![GitHub stars](https://img.shields.io/github/stars/lucap123/Scope?style=flat)
![GitHub issues](https://img.shields.io/github/issues/lucap123/Scope)
![Go version](https://img.shields.io/badge/go-1.22-blue)
![License](https://img.shields.io/badge/license-MIT-green)

<p align="center">
  <video src="scope_demo_v2.mp4" width="100%" controls autoplay loop muted></video>
</p>

> One session. Every tool.

Stop copy-pasting your auth token. **scope** is a lightweight CLI session manager for web hacking set your headers, proxy, and base URL once, then run curl, ffuf, nuclei, sqlmap and more with everything auto-injected.

```bash
scope new bugcrowd-target --url https://api.target.com
scope set header "Authorization: Bearer eyJhbG..."
scope set header "Cookie: session=abc123"
scope set proxy http://127.0.0.1:8080
scope use bugcrowd-target

scope curl /users           # → curl -H "Authorization: ..." https://api.target.com/users
scope ffuf -w words.txt -u /FUZZ
scope nuclei -t exposures/
```

---

## Install

```bash
go install github.com/lucap123/scope@latest
```

Or build from source:

```bash
git clone https://github.com/lucap123/scope
cd scope
make install
```

---

## Usage

### Sessions

```bash
scope new <name>                          # create a session
scope new <name> --url https://target.com

scope set <name> url https://target.com  # set base URL
scope set <name> header "Authorization: Bearer <token>"
scope set <name> header "Cookie: session=abc123"
scope set <name> proxy http://127.0.0.1:8080

scope unset <name> header Authorization  # remove a header

scope use <name>                          # set active session
scope show                                # inspect active session
scope show --reveal                       # show full header values (not masked)
scope list                                # list all sessions
scope delete <name>                       # delete a session
```

Set `SCOPE_SESSION=<name>` in your shell to override the active session per terminal without touching the global state.

### Tool wrapping

```bash
scope curl /users
scope curl -X POST /login -d '{"user":"test","pass":"test"}'

scope ffuf -w wordlist.txt -u /FUZZ
scope ffuf -w wordlist.txt -u /api/FUZZ

scope http GET /users
scope http POST /login user=test pass=test

scope nuclei -t exposures/
scope nuclei -t cves/ -severity high,critical

scope sqlmap -u /search?q=test --dbs
```

Paths starting with `/` are automatically prefixed with the session base URL. All other arguments are passed through unchanged.

Use `--verbose` to see the full command before it runs:

```bash
scope --verbose curl /users
# [scope] curl -H "Authorization: Bearer ey..." --proxy http://127.0.0.1:8080 https://api.target.com/users
```

### Export

```bash
# Generate a copy-paste ready command
scope export curl /users
scope export ffuf -w words.txt -u /FUZZ

# Export the session itself
scope export <name> --format json > session.json
scope export <name> --format env
```

### Import from Burp Suite

Right-click a request in Burp → Save item → import it:

```bash
scope import burp burp-export.xml
```

Parses base URL, auth headers, and cookies. Skips Host, Content-Length, and other per-request headers automatically.

### Notes

```bash
scope note "Admin panel at /manage no auth check"
scope note "v1 endpoint still active at /api/v1/auth"
scope notes
```

---

## How sessions are stored

Plain JSON files in `~/.scope/sessions/`. No cloud, no account, no telemetry.

```
~/.scope/
  sessions/
    bugcrowd-target.json
    hackerone-acme.json
  active
```

Files are `chmod 600`. Tokens are stored in plaintext this is a hacking tool, you control your own security.

Override the storage directory with `SCOPE_HOME=/path/to/dir`.

---

## Build

```bash
make build    # builds ./scope
make install  # installs to $GOPATH/bin
make test     # runs tests
make release  # builds for linux/mac/windows
```

---

## Roadmap

| Version | Features |
|---------|----------|
| v0.1 | Sessions, headers, proxy, curl/ffuf/nuclei/sqlmap wrapping ✓ |
| v0.2 | Burp import, session export, notes ✓ |
| v0.3 | Token auto-refresh via hook script |
| v0.4 | Team sharing with optional encryption |
| v1.0 | Plugin system for community tools |

---

## License

MIT
