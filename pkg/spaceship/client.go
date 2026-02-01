package spaceship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	APIKey     string
	APISecret  string
	HTTPClient *http.Client
}

func NewClient(apiKey, apiSecret string, httpClient *http.Client) *Client {
	return &Client{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		HTTPClient: httpClient,
	}
}

type DNSSaveRequest struct {
	Force bool           `json:"force"`
	Items []DNSTXTRecord `json:"items"`
}

type DNSTXTRecord struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Name  string `json:"name"`
	TTL   int    `json:"ttl,omitempty"`
}

func (c *Client) AddTXTRecord(domain, name, value string, ttl int) error {
	record := DNSTXTRecord{
		Type:  "TXT",
		Name:  name,
		Value: value,
		TTL:   ttl,
	}

	request := DNSSaveRequest{
		Force: true,
		Items: []DNSTXTRecord{record},
	}

	payload, _ := json.Marshal(request)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("https://spaceship.dev/api/v1/dns/records/%s", domain), bytes.NewBuffer(payload))
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("X-API-Secret", c.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to add TXT record: %s", responseDetail(resp))
	}

	return nil
}

func (c *Client) RemoveTXTRecord(domain, name, value string) error {
	records := []DNSTXTRecord{
		{
			Type:  "TXT",
			Name:  name,
			Value: value,
		},
	}

	payload, _ := json.Marshal(records)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("https://spaceship.dev/api/v1/dns/records/%s", domain), bytes.NewBuffer(payload))
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("X-API-Secret", c.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to remove TXT record: %s", responseDetail(resp))
	}

	return nil
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func responseDetail(resp *http.Response) string {
	if resp == nil {
		return "no response"
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.Status
	}
	if len(body) == 0 {
		return resp.Status
	}
	return fmt.Sprintf("%s: %s", resp.Status, string(body))
}
