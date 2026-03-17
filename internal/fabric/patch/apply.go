package patch

import "os"

func Apply(path string, content []byte) error {
	return os.WriteFile(path, content, 0o644)
}
