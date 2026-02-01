package main

import (
	"os"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/neoscrib/cert-manager-spaceship/pkg/spaceship"
)

var (
	getenv           = os.Getenv
	newClient        = spaceship.NewClient
	newSolver        = spaceship.NewSolver
	runWebhookServer = cmd.RunWebhookServer
)

func main() {
	apiKey := getenv("SPACESHIP_API_KEY")
	apiSecret := getenv("SPACESHIP_API_SECRET")
	client := newClient(apiKey, apiSecret, nil)

	solver := newSolver(client)

	runWebhookServer("spaceship-webhook", solver)
}
