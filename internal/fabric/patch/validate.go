package patch

import "fmt"

func ValidatePatch(path string, content []byte) error {
	if path == "" {
		return fmt.Errorf("patch path is required")
	}
	if len(content) == 0 {
		return fmt.Errorf("patch content is required")
	}
	return nil
}
