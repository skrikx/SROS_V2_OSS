package adapters

import (
	"context"
	"fmt"
	"net/http"

	ctools "srosv2/contracts/tools"
)

type LocalHTTP struct{}

func (LocalHTTP) Name() string { return "local_http" }

func (LocalHTTP) Invoke(_ context.Context, req map[string]any, env ctools.AuthEnvelope) (map[string]any, error) {
	url, _ := req["url"].(string)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}
	return map[string]any{
		"connector": env.Connector,
		"url":       url,
		"status":    resp.StatusCode,
	}, nil
}
