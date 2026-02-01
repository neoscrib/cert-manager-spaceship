package spaceship

import (
	"testing"

	cmacme "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	restclient "k8s.io/client-go/rest"
)

type solverClientStub struct {
	addArgs    []string
	removeArgs []string
	addErr     error
	removeErr  error
	addTTL     int
}

func (s *solverClientStub) AddTXTRecord(domain, name, value string, ttl int) error {
	s.addArgs = []string{domain, name, value}
	s.addTTL = ttl
	return s.addErr
}

func (s *solverClientStub) RemoveTXTRecord(domain, name, value string) error {
	s.removeArgs = []string{domain, name, value}
	return s.removeErr
}

func TestSolverPresent(t *testing.T) {
	stub := &solverClientStub{}
	solver := &Solver{client: stub}

	ch := &cmacme.ChallengeRequest{
		ResolvedFQDN: "_acme-challenge.example.com.",
		ResolvedZone: "example.com.",
		Key:          "token",
	}

	if err := solver.Present(ch); err != nil {
		t.Fatalf("Present error: %v", err)
	}
	if len(stub.addArgs) != 3 {
		t.Fatalf("addArgs len = %d", len(stub.addArgs))
	}
	if stub.addArgs[0] != "example.com." || stub.addArgs[1] != "_acme-challenge.example.com." || stub.addArgs[2] != "token" {
		t.Fatalf("addArgs = %#v", stub.addArgs)
	}
	if stub.addTTL != 60 {
		t.Fatalf("addTTL = %d", stub.addTTL)
	}
}

func TestSolverCleanUp(t *testing.T) {
	stub := &solverClientStub{}
	solver := &Solver{client: stub}

	ch := &cmacme.ChallengeRequest{
		ResolvedFQDN: "_acme-challenge.example.com.",
		ResolvedZone: "example.com.",
		Key:          "token",
	}

	if err := solver.CleanUp(ch); err != nil {
		t.Fatalf("CleanUp error: %v", err)
	}
	if len(stub.removeArgs) != 3 {
		t.Fatalf("removeArgs len = %d", len(stub.removeArgs))
	}
	if stub.removeArgs[0] != "example.com." || stub.removeArgs[1] != "_acme-challenge.example.com." || stub.removeArgs[2] != "token" {
		t.Fatalf("removeArgs = %#v", stub.removeArgs)
	}
}

func TestSolverInitializeNoop(t *testing.T) {
	solver := &Solver{}
	if err := solver.Initialize(&restclient.Config{}, make(chan struct{})); err != nil {
		t.Fatalf("Initialize error: %v", err)
	}
}
