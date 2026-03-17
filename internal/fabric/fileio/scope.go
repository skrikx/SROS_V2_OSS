package fileio

import (
	"fmt"
	"path/filepath"
	"strings"
)

func EnsureWithinRoot(root, target string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absTarget, absRoot) {
		return fmt.Errorf("target %s escapes root %s", absTarget, absRoot)
	}
	return nil
}
