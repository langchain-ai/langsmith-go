//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

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
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return &oauthRefreshLock{file: file}, nil
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) && !errors.Is(err, syscall.EAGAIN) {
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
	err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	closeErr := l.file.Close()
	if err != nil {
		return fmt.Errorf("release OAuth refresh lock: %w", err)
	}
	return closeErr
}
