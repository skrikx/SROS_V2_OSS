package mcpclient

import (
	"fmt"

	"srosv2/internal/fabric/registry"
)

type Client struct {
	registry *registry.Registry
}

func New(reg *registry.Registry) *Client {
	return &Client{registry: reg}
}

func (c *Client) Ingest(path string) (map[string]any, error) {
	manifest, err := Ingest(path, c.registry)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"accepted": true,
		"summary":  fmt.Sprintf("MCP capability %s normalized and admitted", manifest.Name),
		"manifest": manifest,
	}, nil
}
