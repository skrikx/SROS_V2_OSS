package web

import "strings"

func ExtractTitle(body string) string {
	start := strings.Index(strings.ToLower(body), "<title>")
	end := strings.Index(strings.ToLower(body), "</title>")
	if start == -1 || end == -1 || end <= start+7 {
		return ""
	}
	return strings.TrimSpace(body[start+7 : end])
}
