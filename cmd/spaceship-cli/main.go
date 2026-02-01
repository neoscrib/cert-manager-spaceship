package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/neoscrib/cert-manager-spaceship/pkg/spaceship"
)

const defaultTimeout = 10 * time.Second

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "add":
		runAdd(os.Args[2:])
	case "remove":
		runRemove(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func runAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	domain := fs.String("domain", "", "domain to update (required)")
	name := fs.String("name", "", "record name (required)")
	value := fs.String("value", "", "record value (required)")
	ttl := fs.Int("ttl", 60, "record TTL seconds")
	apiKey := fs.String("api-key", os.Getenv("SPACESHIP_API_KEY"), "Spaceship API key (or SPACESHIP_API_KEY)")
	apiSecret := fs.String("api-secret", os.Getenv("SPACESHIP_API_SECRET"), "Spaceship API secret (or SPACESHIP_API_SECRET)")
	timeout := fs.Duration("timeout", defaultTimeout, "HTTP timeout")
	fs.Parse(args)

	if !validateRequired(*domain, *name, *value, *apiKey, *apiSecret) {
		fs.Usage()
		os.Exit(2)
	}

	client := spaceship.NewClient(*apiKey, *apiSecret, &http.Client{Timeout: *timeout})
	if err := client.AddTXTRecord(*domain, *name, *value, *ttl); err != nil {
		fmt.Fprintf(os.Stderr, "add failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")
}

func runRemove(args []string) {
	fs := flag.NewFlagSet("remove", flag.ExitOnError)
	domain := fs.String("domain", "", "domain to update (required)")
	name := fs.String("name", "", "record name (required)")
	value := fs.String("value", "", "record value (required)")
	apiKey := fs.String("api-key", os.Getenv("SPACESHIP_API_KEY"), "Spaceship API key (or SPACESHIP_API_KEY)")
	apiSecret := fs.String("api-secret", os.Getenv("SPACESHIP_API_SECRET"), "Spaceship API secret (or SPACESHIP_API_SECRET)")
	timeout := fs.Duration("timeout", defaultTimeout, "HTTP timeout")
	fs.Parse(args)

	if !validateRequired(*domain, *name, *value, *apiKey, *apiSecret) {
		fs.Usage()
		os.Exit(2)
	}

	client := spaceship.NewClient(*apiKey, *apiSecret, &http.Client{Timeout: *timeout})
	if err := client.RemoveTXTRecord(*domain, *name, *value); err != nil {
		fmt.Fprintf(os.Stderr, "remove failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")
}

func validateRequired(domain, name, value, apiKey, apiSecret string) bool {
	if domain == "" || name == "" || value == "" || apiKey == "" || apiSecret == "" {
		return false
	}
	return true
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  spaceship-cli add --domain example.com --name _acme-challenge --value token [--ttl 60]")
	fmt.Fprintln(os.Stderr, "  spaceship-cli remove --domain example.com --name _acme-challenge --value token")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Environment:")
	fmt.Fprintln(os.Stderr, "  SPACESHIP_API_KEY, SPACESHIP_API_SECRET")
}
