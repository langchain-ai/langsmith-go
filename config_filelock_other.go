//go:build !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd

package langsmith

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type oauthRefreshLock struct {
	dir *oauthRefreshDirLock
}

func acquireOAuthRefreshLock(ctx context.Context, path string) (*oauthRefreshLock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("create OAuth refresh lock directory: %w", err)
	}

	lock, err := acquireOAuthRefreshDirLock(ctx, path+".lock", time.Now)
	if err != nil {
		return nil, err
	}
	return &oauthRefreshLock{dir: lock}, nil
}

func (l *oauthRefreshLock) Unlock() error {
	return l.dir.unlock()
}
