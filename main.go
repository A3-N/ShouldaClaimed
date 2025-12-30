package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/A3-N/ShouldaClaimed/pkg/cli"
	"github.com/A3-N/ShouldaClaimed/pkg/generator"
)

const banner = `
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
`

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	fmt.Print(banner)
	cli.PrintInfo("ShouldaClaimed CLI initialized")

	switch os.Args[1] {
	case "create":
		handleCreate(os.Args[2:])
	case "publish":
		handlePublish(os.Args[2:])
	case "-h", "--help":
		printUsage()
	default:
		cli.PrintError("Unknown subcommand: %s", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(banner)
	fmt.Println("Usage: ShouldaClaimed <command> [flags]")
	fmt.Println("\nCommands:")
	fmt.Println("  create   Generate a package without publishing")
	fmt.Println("  publish  Generate and publish a package")
	fmt.Println("\nRun 'ShouldaClaimed <command> -h' for more information.")
}

func handleCreate(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	cfg, err := cli.ParseFlags(fs, args)
	if err != nil {
		cli.PrintError("Error parsing flags: %v", err)
		fs.PrintDefaults()
		os.Exit(1)
	}

	runChecks(cfg)

	if err := generator.GeneratePackage(cfg); err != nil {
		cli.PrintError("Generation failed: %v", err)
		os.Exit(1)
	}
	cli.PrintSuccess("Package generated successfully at %s/%s", cfg.OutDir, cfg.Target)
}

func handlePublish(args []string) {
	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	cfg, err := cli.ParseFlags(fs, args)
	if err != nil {
		cli.PrintError("Error parsing flags: %v", err)
		fs.PrintDefaults()
		os.Exit(1)
	}

	runChecks(cfg)

	// 1. Authenticate
	cli.PrintInfo("Ensuring NPM authentication...")
	if err := ensureNPMAuth(); err != nil {
		cli.PrintError("Authentication failed: %v", err)
		os.Exit(1)
	}

	// 2. Generate
	if err := generator.GeneratePackage(cfg); err != nil {
		cli.PrintError("Generation failed: %v", err)
		os.Exit(1)
	}

	// 3. Publish
	pkgDir := filepath.Join(cfg.OutDir, cfg.Target)
	cli.PrintInfo("Publishing package from %s...", pkgDir)

	cmd := exec.Command("npm", "publish")
	cmd.Dir = pkgDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		cli.PrintError("Publishing failed: %v", err)
		os.Exit(1)
	}
	cli.PrintSuccess("Package published successfully!")
}

func ensureNPMAuth() error {
	// Check whoami
	whoami := exec.Command("npm", "whoami")
	if err := whoami.Run(); err == nil {
		// Already logged in
		return nil
	}

	cli.PrintInfo("Not logged in. Initiating npm login...")
	// Attempt login with --auth-type=web
	login := exec.Command("npm", "login", "--auth-type=web")
	login.Stdin = os.Stdin
	login.Stdout = os.Stdout
	login.Stderr = os.Stderr

	if err := login.Run(); err != nil {
		return fmt.Errorf("npm login failed")
	}
	return nil
}

func runChecks(cfg *cli.Config) {
	if cfg.Force {
		cli.PrintInfo("Force mode enabled. Skipping pre-flight checks.")
		return
	}

	if err := checkNPMInstalled(); err != nil {
		cli.PrintError("Pre-flight check failed: %v", err)
		os.Exit(1)
	}

	// Logic:
	// 1. Check if local directory exists. If YES -> Skip exist check (assume update).
	// 2. If NO -> Check npm view. If 0 (exists in registry) -> Fail.

	localPath := filepath.Join(cfg.OutDir, cfg.Target)
	if _, err := os.Stat(localPath); !os.IsNotExist(err) {
		// Local directory exists.
		cli.PrintInfo("Local directory %s exists. Skipping registry check.", localPath)
		return
	}

	if err := checkPackageExists(cfg.Target); err != nil {
		cli.PrintError("Pre-flight check failed: %v", err)
		os.Exit(1)
	}
}

func checkNPMInstalled() error {
	cmd := exec.Command("npm", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm is not installed or not in PATH")
	}
	return nil
}

func checkPackageExists(pkgName string) error {
	// npm view <pkg> returns 0 if exists, 1 if 404
	cmd := exec.Command("npm", "view", pkgName)
	if err := cmd.Run(); err == nil {
		// Success means package exists -> Error in this context
		return fmt.Errorf("package '%s' already exists in the registry", pkgName)
	}
	// Error means 404 (not found) -> Good
	return nil
}
