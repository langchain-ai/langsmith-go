//go:build !windows

package langsmith

import "os"

func replaceFile(from, to string) error {
	return os.Rename(from, to)
}
