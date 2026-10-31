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

// The OAuth refresh lock is shared by every LangSmith SDK (Python, JS/TS, Go,
// Java) that reads the same profile config: an exclusive directory created with
// an atomic mkdir at "<config>.oauth.lock.lock". Its "created_at" file holds a
// UTC timestamp and a random owner ID. A lock older than
// oauthRefreshLockStaleAfter is treated as abandoned (its holder crashed) and
// taken over. Holders bound their whole refresh to tokenRefreshTimeout, well
// under that threshold, so a live holder is never evicted.
const (
	oauthRefreshLockPollInterval  = 10 * time.Millisecond
	oauthRefreshLockStaleAfter    = 15 * time.Second
	oauthRefreshLockTimestampFile = "created_at"
	// On Windows a file cannot be deleted while another process has it open
	// without delete sharing (Go and Python open files that way), so removing
	// the lock fails if a waiter is reading created_at at that moment. Reads
	// take microseconds; a few short retries absorb the overlap.
	oauthRefreshLockRemoveAttempts = 10
)

// oauthRefreshLockRemoveAll is os.RemoveAll; tests replace it to simulate
// Windows sharing violations.
var oauthRefreshLockRemoveAll = os.RemoveAll

type oauthRefreshLock struct {
	path  string
	owner string
}

// acquireOAuthRefreshLock takes the shared refresh lock for the config whose
// lock path is path ("<config>.oauth.lock"); the lock directory is path+".lock".
func acquireOAuthRefreshLock(ctx context.Context, path string) (*oauthRefreshLock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("create OAuth refresh lock directory: %w", err)
	}
	return acquireOAuthRefreshDirLock(ctx, path+".lock", time.Now)
}

func acquireOAuthRefreshDirLock(ctx context.Context, lockDir string, now func() time.Time) (*oauthRefreshLock, error) {
	owner := newOAuthRefreshLockOwner()
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := os.Mkdir(lockDir, 0700); err == nil {
			if err := writeOAuthRefreshLockMetadata(lockDir, now(), owner); err != nil {
				_ = removeOAuthRefreshLockDir(lockDir)
				return nil, fmt.Errorf("write OAuth refresh lock metadata: %w", err)
			}
			return &oauthRefreshLock{path: lockDir, owner: owner}, nil
		} else if !os.IsExist(err) {
			return nil, fmt.Errorf("acquire OAuth refresh lock: %w", err)
		}

		if removed, err := removeStaleOAuthRefreshLock(lockDir, now()); err != nil {
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

// Unlock removes the lock directory if this holder still owns it.
func (l *oauthRefreshLock) Unlock() error {
	if owner, ok := oauthRefreshLockOwner(l.path); !ok || owner != l.owner {
		return nil
	}
	return removeOAuthRefreshLockDir(l.path)
}

func writeOAuthRefreshLockMetadata(lockDir string, now time.Time, owner string) error {
	return os.WriteFile(
		filepath.Join(lockDir, oauthRefreshLockTimestampFile),
		[]byte(now.UTC().Format(time.RFC3339Nano)+"\n"+owner+"\n"),
		0600,
	)
}

func removeStaleOAuthRefreshLock(lockDir string, now time.Time) (bool, error) {
	createdAt, ok := oauthRefreshLockCreatedAt(lockDir)
	if !ok || now.Sub(createdAt) <= oauthRefreshLockStaleAfter {
		return false, nil
	}
	if err := removeOAuthRefreshLockDir(lockDir); err != nil {
		return false, fmt.Errorf("remove stale OAuth refresh lock: %w", err)
	}
	return true, nil
}

// removeOAuthRefreshLockDir removes the lock directory, retrying briefly when
// another process holds one of its files open.
func removeOAuthRefreshLockDir(lockDir string) error {
	var err error
	for attempt := range oauthRefreshLockRemoveAttempts {
		if attempt > 0 {
			time.Sleep(oauthRefreshLockPollInterval)
		}
		if err = oauthRefreshLockRemoveAll(lockDir); err == nil {
			return nil
		}
	}
	return err
}

// oauthRefreshLockCreatedAt reads the lock's timestamp, falling back to the
// directory's mtime when the metadata is missing or malformed (e.g. the holder
// crashed between mkdir and writing it).
func oauthRefreshLockCreatedAt(lockDir string) (time.Time, bool) {
	lines, ok := oauthRefreshLockMetadata(lockDir)
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

func oauthRefreshLockOwner(lockDir string) (string, bool) {
	lines, ok := oauthRefreshLockMetadata(lockDir)
	if !ok || len(lines) < 2 {
		return "", false
	}
	owner := strings.TrimSpace(lines[1])
	return owner, owner != ""
}

func oauthRefreshLockMetadata(lockDir string) ([]string, bool) {
	data, err := os.ReadFile(filepath.Join(lockDir, oauthRefreshLockTimestampFile))
	if err != nil {
		return nil, false
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n"), true
}

func newOAuthRefreshLockOwner() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return hex.EncodeToString(b[:])
	}
	return fmt.Sprintf("%d-%d", os.Getpid(), time.Now().UnixNano())
}
