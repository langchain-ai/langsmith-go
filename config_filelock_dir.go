package langsmith

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	oauthRefreshLockPollInterval  = 10 * time.Millisecond
	oauthRefreshLockStaleAfter    = 10 * time.Second
	oauthRefreshLockTimestampFile = "created_at"
)

type oauthRefreshDirLock struct {
	path  string
	owner string
}

func acquireOAuthRefreshDirLock(ctx context.Context, lockDir string, now func() time.Time) (*oauthRefreshDirLock, error) {
	owner := newOAuthRefreshDirLockOwner()
	for {
		if err := os.Mkdir(lockDir, 0700); err == nil {
			if err := writeOAuthRefreshLockMetadata(lockDir, now(), owner); err != nil {
				_ = os.RemoveAll(lockDir)
				return nil, fmt.Errorf("write OAuth refresh lock metadata: %w", err)
			}
			return &oauthRefreshDirLock{path: lockDir, owner: owner}, nil
		} else if !os.IsExist(err) {
			return nil, fmt.Errorf("acquire OAuth refresh lock: %w", err)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if removed, err := removeStaleOAuthRefreshDirLock(lockDir, now()); err != nil {
			return nil, err
		} else if removed {
			continue
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(oauthRefreshLockPollInterval):
		}
	}
}

func (l *oauthRefreshDirLock) unlock() error {
	if owner, ok := oauthRefreshDirLockOwner(l.path); !ok || owner != l.owner {
		return nil
	}
	return os.RemoveAll(l.path)
}

func writeOAuthRefreshLockMetadata(lockDir string, now time.Time, owner string) error {
	return os.WriteFile(
		filepath.Join(lockDir, oauthRefreshLockTimestampFile),
		[]byte(now.UTC().Format(time.RFC3339Nano)+"\n"+owner+"\n"),
		0600,
	)
}

func removeStaleOAuthRefreshDirLock(lockDir string, now time.Time) (bool, error) {
	createdAt, ok := oauthRefreshDirLockCreatedAt(lockDir)
	if !ok || now.Sub(createdAt) <= oauthRefreshLockStaleAfter {
		return false, nil
	}
	if err := os.RemoveAll(lockDir); err != nil {
		return false, fmt.Errorf("remove stale OAuth refresh lock: %w", err)
	}
	return true, nil
}

func oauthRefreshDirLockCreatedAt(lockDir string) (time.Time, bool) {
	lines, ok := oauthRefreshDirLockMetadata(lockDir)
	if ok && len(lines) > 0 {
		createdAt, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(lines[0]))
		if err == nil {
			return createdAt, true
		}
	}
	info, err := os.Stat(lockDir)
	if err != nil {
		return time.Time{}, false
	}
	return info.ModTime(), true
}

func oauthRefreshDirLockOwner(lockDir string) (string, bool) {
	lines, ok := oauthRefreshDirLockMetadata(lockDir)
	if !ok || len(lines) < 2 {
		return "", false
	}
	owner := strings.TrimSpace(lines[1])
	return owner, owner != ""
}

func oauthRefreshDirLockMetadata(lockDir string) ([]string, bool) {
	data, err := os.ReadFile(filepath.Join(lockDir, oauthRefreshLockTimestampFile))
	if err != nil {
		return nil, false
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n"), true
}

func newOAuthRefreshDirLockOwner() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return hex.EncodeToString(b[:])
	}
	return fmt.Sprintf("%d-%d", os.Getpid(), time.Now().UnixNano())
}
