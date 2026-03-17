package patch

import "strings"

func UnifiedPreview(before, after string) string {
	return "--- before\n+++ after\n-" + strings.ReplaceAll(before, "\n", "\n-") + "\n+" + strings.ReplaceAll(after, "\n", "\n+")
}
