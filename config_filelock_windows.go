//go:build windows

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows"
)

// oauthRefreshLockAllBytes locks the whole file, including past its end, so
// the lock works on the empty lock file.
const oauthRefreshLockAllBytes = ^uint32(0)

// oauthRefreshLock is a LockFileEx lock on the lock file. Windows releases it
// when the handle closes, including when the process dies, and nothing has to
// be deleted to release it: Go opens files without FILE_SHARE_DELETE, so a
// lock released by deleting a file fails while any waiter is reading it.
type oauthRefreshLock struct {
	file *os.File
}

func acquireOAuthRefreshLock(ctx context.Context, path string) (*oauthRefreshLock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("create OAuth refresh lock directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("open OAuth refresh lock file: %w", err)
	}

	for {
		err = windows.LockFileEx(
			windows.Handle(file.Fd()),
			windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
			0,
			oauthRefreshLockAllBytes,
			oauthRefreshLockAllBytes,
			new(windows.Overlapped),
		)
		if err == nil {
			return &oauthRefreshLock{file: file}, nil
		}
		if !errors.Is(err, windows.ERROR_LOCK_VIOLATION) {
			_ = file.Close()
			return nil, fmt.Errorf("acquire OAuth refresh lock: %w", err)
		}

		select {
		case <-ctx.Done():
			_ = file.Close()
			return nil, ctx.Err()
		case <-time.After(oauthRefreshLockPollInterval):
		}
	}
}

func (l *oauthRefreshLock) Unlock() error {
	err := windows.UnlockFileEx(
		windows.Handle(l.file.Fd()),
		0,
		oauthRefreshLockAllBytes,
		oauthRefreshLockAllBytes,
		new(windows.Overlapped),
	)
	closeErr := l.file.Close()
	if err != nil {
		return fmt.Errorf("release OAuth refresh lock: %w", err)
	}
	return closeErr
}
