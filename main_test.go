package main

import (
	"net/http"
	"testing"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook"
	"github.com/neoscrib/cert-manager-spaceship/pkg/spaceship"
)

func TestMainWiresDependencies(t *testing.T) {
	origGetenv := getenv
	origNewClient := newClient
	origNewSolver := newSolver
	origRun := runWebhookServer

	t.Cleanup(func() {
		getenv = origGetenv
		newClient = origNewClient
		newSolver = origNewSolver
		runWebhookServer = origRun
	})

	getenv = func(key string) string {
		switch key {
		case "SPACESHIP_API_KEY":
			return "key"
		case "SPACESHIP_API_SECRET":
			return "secret"
		default:
			return ""
		}
	}

	var gotKey string
	var gotSecret string
	newClient = func(apiKey, apiSecret string, _ *http.Client) *spaceship.Client {
		gotKey = apiKey
		gotSecret = apiSecret
		return &spaceship.Client{APIKey: apiKey, APISecret: apiSecret}
	}

	newSolver = func(client spaceship.SolverClient) *spaceship.Solver {
		return spaceship.NewSolver(client)
	}

	var gotGroup string
	var gotHooks []webhook.Solver
	runWebhookServer = func(groupName string, hooks ...webhook.Solver) {
		gotGroup = groupName
		gotHooks = hooks
	}

	main()

	if gotKey != "key" || gotSecret != "secret" {
		t.Fatalf("client args = %q/%q", gotKey, gotSecret)
	}
	if gotGroup != "spaceship-webhook" {
		t.Fatalf("group = %q", gotGroup)
	}
	if len(gotHooks) != 1 || gotHooks[0] == nil {
		t.Fatalf("hooks = %#v", gotHooks)
	}
}
