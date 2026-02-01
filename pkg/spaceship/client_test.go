package spaceship

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func okResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}
}

func statusResponse(code int, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestAddTXTRecordSuccess(t *testing.T) {
	client := &Client{
		APIKey:    "key",
		APISecret: "secret",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodPut {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodPut)
				}
				if r.URL.String() != "https://spaceship.dev/api/v1/dns/records/example.com" {
					t.Fatalf("url = %s", r.URL.String())
				}
				if got := r.Header.Get("X-API-Key"); got != "key" {
					t.Fatalf("X-API-Key = %s", got)
				}
				if got := r.Header.Get("X-API-Secret"); got != "secret" {
					t.Fatalf("X-API-Secret = %s", got)
				}
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Fatalf("Content-Type = %s", got)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}
				var req DNSSaveRequest
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if !req.Force {
					t.Fatalf("Force = false")
				}
				if len(req.Items) != 1 {
					t.Fatalf("Items len = %d", len(req.Items))
				}
				item := req.Items[0]
				if item.Type != "TXT" || item.Name != "name" || item.Value != "value" || item.TTL != 60 {
					t.Fatalf("item = %+v", item)
				}

				return okResponse(), nil
			}),
		},
	}

	if err := client.AddTXTRecord("example.com", "name", "value", 60); err != nil {
		t.Fatalf("AddTXTRecord error: %v", err)
	}
}

func TestAddTXTRecordStatusError(t *testing.T) {
	client := &Client{
		APIKey:    "key",
		APISecret: "secret",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return statusResponse(http.StatusBadRequest, "bad request"), nil
			}),
		},
	}

	if err := client.AddTXTRecord("example.com", "name", "value", 60); err == nil {
		t.Fatalf("expected error")
	} else if !strings.Contains(err.Error(), "400 Bad Request") {
		t.Fatalf("error = %v", err)
	}
}

func TestAddTXTRecordNoContentSuccess(t *testing.T) {
	client := &Client{
		APIKey:    "key",
		APISecret: "secret",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return statusResponse(http.StatusNoContent, ""), nil
			}),
		},
	}

	if err := client.AddTXTRecord("example.com", "name", "value", 60); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestAddTXTRecordTransportError(t *testing.T) {
	wantErr := errors.New("transport error")

	client := &Client{
		APIKey:    "key",
		APISecret: "secret",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return nil, wantErr
			}),
		},
	}

	if err := client.AddTXTRecord("example.com", "name", "value", 60); !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestRemoveTXTRecordSuccess(t *testing.T) {
	client := &Client{
		APIKey:    "key",
		APISecret: "secret",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodDelete {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodDelete)
				}
				if r.URL.String() != "https://spaceship.dev/api/v1/dns/records/example.com" {
					t.Fatalf("url = %s", r.URL.String())
				}
				if got := r.Header.Get("X-API-Key"); got != "key" {
					t.Fatalf("X-API-Key = %s", got)
				}
				if got := r.Header.Get("X-API-Secret"); got != "secret" {
					t.Fatalf("X-API-Secret = %s", got)
				}
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Fatalf("Content-Type = %s", got)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}
				var records []DNSTXTRecord
				if err := json.Unmarshal(body, &records); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if len(records) != 1 {
					t.Fatalf("records len = %d", len(records))
				}
				record := records[0]
				if record.Type != "TXT" || record.Name != "name" || record.Value != "value" {
					t.Fatalf("record = %+v", record)
				}

				return okResponse(), nil
			}),
		},
	}

	if err := client.RemoveTXTRecord("example.com", "name", "value"); err != nil {
		t.Fatalf("RemoveTXTRecord error: %v", err)
	}
}

func TestRemoveTXTRecordStatusError(t *testing.T) {
	client := &Client{
		APIKey:    "key",
		APISecret: "secret",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return statusResponse(http.StatusBadRequest, "bad request"), nil
			}),
		},
	}

	if err := client.RemoveTXTRecord("example.com", "name", "value"); err == nil {
		t.Fatalf("expected error")
	} else if !strings.Contains(err.Error(), "400 Bad Request") {
		t.Fatalf("error = %v", err)
	}
}

func TestRemoveTXTRecordNoContentSuccess(t *testing.T) {
	client := &Client{
		APIKey:    "key",
		APISecret: "secret",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return statusResponse(http.StatusNoContent, ""), nil
			}),
		},
	}

	if err := client.RemoveTXTRecord("example.com", "name", "value"); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestRemoveTXTRecordTransportError(t *testing.T) {
	wantErr := errors.New("transport error")

	client := &Client{
		APIKey:    "key",
		APISecret: "secret",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return nil, wantErr
			}),
		},
	}

	if err := client.RemoveTXTRecord("example.com", "name", "value"); !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
