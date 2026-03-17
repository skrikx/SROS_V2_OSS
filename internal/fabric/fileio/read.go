package fileio

import "os"

func Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}
