package spaceship

import (
	restclient "k8s.io/client-go/rest"

	cmacme "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
)

const SolverName = "spaceship"

type SolverClient interface {
	AddTXTRecord(domain, name, value string, ttl int) error
	RemoveTXTRecord(domain, name, value string) error
}

type Solver struct {
	client SolverClient
}

func NewSolver(client SolverClient) *Solver {
	return &Solver{
		client: client,
	}
}

func (s *Solver) Name() string {
	return SolverName
}

func (s *Solver) Present(ch *cmacme.ChallengeRequest) error {
	fqdn := ch.ResolvedFQDN
	value := ch.Key
	domain := ch.ResolvedZone

	return s.client.AddTXTRecord(domain, fqdn, value, 60)
}

func (s *Solver) CleanUp(ch *cmacme.ChallengeRequest) error {
	fqdn := ch.ResolvedFQDN
	value := ch.Key
	domain := ch.ResolvedZone

	return s.client.RemoveTXTRecord(domain, fqdn, value)
}

func (s *Solver) Initialize(_ *restclient.Config, _ <-chan struct{}) error {
	return nil
}
