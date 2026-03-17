package shell

import "strings"

func SafeCommand(name string) bool {
	blocked := []string{"powershell -enc", "curl |", "wget |"}
	candidate := strings.ToLower(strings.TrimSpace(name))
	for _, item := range blocked {
		if strings.Contains(candidate, item) {
			return false
		}
	}
	return candidate != ""
}
