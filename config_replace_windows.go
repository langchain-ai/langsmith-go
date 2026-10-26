//go:build windows

package langsmith

import (
	"errors"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// replaceFileTimeout bounds how long replaceFile waits for readers to close
// the destination.
const replaceFileTimeout = 2 * time.Second

// replaceFile renames from over to. Windows refuses to replace a file that
// another handle has open without FILE_SHARE_DELETE, and Go opens files without
// it, so a command reading the config at the moment a refresh saves it would
// fail the save. Readers close the file quickly, so retry until they do.
func replaceFile(from, to string) error {
	deadline := time.Now().Add(replaceFileTimeout)
	wait := time.Millisecond
	for {
		err := os.Rename(from, to)
		if err == nil || !replaceFileMayRetry(err) || time.Now().Add(wait).After(deadline) {
			return err
		}
		time.Sleep(wait)
		wait = min(2*wait, 50*time.Millisecond)
	}
}

func replaceFileMayRetry(err error) bool {
	return errors.Is(err, windows.ERROR_ACCESS_DENIED) || errors.Is(err, windows.ERROR_SHARING_VIOLATION)
}
