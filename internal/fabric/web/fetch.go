package web

import (
	"fmt"
	"io"
	"net/http"
)

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("fetch failed with status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
