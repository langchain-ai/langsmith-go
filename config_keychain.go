package langsmith

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"
)

const (
	// secretStoreKeychain is the profile secret_store value saying the
	// profile's API key and OAuth tokens live in the OS keychain (macOS
	// Keychain, Windows Credential Manager, or the Linux Secret Service)
	// rather than in the config file.
	secretStoreKeychain = "keychain"
	keychainService     = "langsmith"
	// keychainTimeout bounds a keychain call, which can hang when the Secret
	// Service has no session bus to answer it.
	keychainTimeout = 3 * time.Second
)

var errKeychainTimeout = errors.New("timed out waiting for the OS keychain")

// keyringSet is keyring.Set, replaced in tests that need writes to fail while
// reads succeed.
var keyringSet = keyring.Set

// profileSecrets is the keychain entry of a profile whose secret_store is
// "keychain", stored as JSON under service "langsmith".
type profileSecrets struct {
	APIKey       string `json:"api_key,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// keychainAccount names a profile's keychain entry: the profile name for the
// default config file, and the profile name qualified by the config path for
// any other, so profiles of the same name in two config files stay apart.
func keychainAccount(path, profileName string) string {
	canonical := canonicalConfigPath(path)
	if home, err := os.UserHomeDir(); err == nil && canonical == canonicalConfigPath(filepath.Join(home, ".langsmith", "config.json")) {
		return profileName
	}
	return profileName + "@" + canonical
}

// canonicalConfigPath makes path absolute and resolves symlinks, so every
// alias of one config file names the same keychain entry. A file that does not
// exist yet resolves through its directory, so the name is the same before and
// after the file is first written.
func canonicalConfigPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		return resolved
	}
	if dir, err := filepath.EvalSymlinks(filepath.Dir(abs)); err == nil {
		return filepath.Join(dir, filepath.Base(abs))
	}
	return abs
}

// withKeychainTimeout runs a keychain call, giving up after keychainTimeout.
// An abandoned call finishes in the background; its result is discarded.
func withKeychainTimeout[T any](fn func() (T, error)) (T, error) {
	type result struct {
		value T
		err   error
	}
	ch := make(chan result, 1)
	go func() {
		value, err := fn()
		ch <- result{value, err}
	}()
	timer := time.NewTimer(keychainTimeout)
	defer timer.Stop()
	select {
	case r := <-ch:
		return r.value, r.err
	case <-timer.C:
		var zero T
		return zero, errKeychainTimeout
	}
}

func readProfileSecrets(path, profileName string) (profileSecrets, error) {
	account := keychainAccount(path, profileName)
	data, err := withKeychainTimeout(func() (string, error) {
		return keyring.Get(keychainService, account)
	})
	if errors.Is(err, keyring.ErrNotFound) {
		return profileSecrets{}, fmt.Errorf("LangSmith profile %q stores its credentials in the OS keychain, but no entry was found; run 'langsmith auth login --profile %s'", profileName, profileName)
	}
	if err != nil {
		return profileSecrets{}, fmt.Errorf("reading LangSmith profile %q credentials from the OS keychain: %w", profileName, err)
	}
	var secrets profileSecrets
	if err := json.Unmarshal([]byte(data), &secrets); err != nil {
		return profileSecrets{}, fmt.Errorf("parsing LangSmith profile %q credentials from the OS keychain: %w", profileName, err)
	}
	return secrets, nil
}

func writeProfileSecrets(path, profileName string, secrets profileSecrets) error {
	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	account := keychainAccount(path, profileName)
	_, err = withKeychainTimeout(func() (struct{}, error) {
		return struct{}{}, keyringSet(keychainService, account, string(data))
	})
	return err
}

// updateProfileSecrets applies update to a profile's keychain entry, starting
// from an empty entry when there is none.
func updateProfileSecrets(path, profileName string, update func(*profileSecrets)) error {
	account := keychainAccount(path, profileName)
	data, err := withKeychainTimeout(func() (string, error) {
		return keyring.Get(keychainService, account)
	})
	var secrets profileSecrets
	switch {
	case errors.Is(err, keyring.ErrNotFound):
	case err != nil:
		return err
	default:
		if err := json.Unmarshal([]byte(data), &secrets); err != nil {
			return err
		}
	}
	update(&secrets)
	return writeProfileSecrets(path, profileName, secrets)
}

func deleteProfileSecrets(path, profileName string) {
	account := keychainAccount(path, profileName)
	_, _ = withKeychainTimeout(func() (struct{}, error) {
		return struct{}{}, keyring.Delete(keychainService, account)
	})
}

// fillProfileSecrets copies a keychain-backed profile's secrets into p. A
// secret the config file itself holds is kept, so the precedence of api_key
// over api_key_env and api_key_cmd is the same wherever the key is stored.
func fillProfileSecrets(path, profileName string, p *configProfile) error {
	if p.SecretStore != secretStoreKeychain {
		return nil
	}
	secrets, err := readProfileSecrets(path, profileName)
	if err != nil {
		return err
	}
	if p.APIKey == "" {
		p.APIKey = secrets.APIKey
	}
	if p.OAuth.AccessToken == "" {
		p.OAuth.AccessToken = secrets.AccessToken
	}
	if p.OAuth.RefreshToken == "" {
		p.OAuth.RefreshToken = secrets.RefreshToken
	}
	return nil
}
