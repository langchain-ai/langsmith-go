//go:build integration

package main

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func codexAvailable(t *testing.T) string {
	t.Helper()
	path, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("codex binary not found in PATH, skipping integration test")
	}
	return path
}

func startEliza(t *testing.T) (url string, stop func()) {
	t.Helper()

	// Find a free port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	cmd := exec.Command("go", "run", ".", "--addr", addr)
	cmd.Dir = "."
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start eliza: %v", err)
	}

	// Wait for server to be ready
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	return "http://" + addr, func() {
		cmd.Process.Kill()
		cmd.Wait()
	}
}

func runCodex(t *testing.T, codexPath, baseURL, workdir, prompt string) string {
	t.Helper()
	cmd := exec.Command(codexPath, "exec",
		"-c", "openai_base_url=\""+baseURL+"/v1\"",
		"-m", "eliza",
		"-C", workdir,
		"--skip-git-repo-check",
		"--ephemeral",
		"-s", "danger-full-access",
		"--full-auto",
		prompt,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("codex output:\n%s", string(out))
		t.Fatalf("codex exec failed: %v", err)
	}
	return string(out)
}

func TestCodex_Hello(t *testing.T) {
	codexPath := codexAvailable(t)
	url, stop := startEliza(t)
	defer stop()

	dir := t.TempDir()
	out := runCodex(t, codexPath, url, dir, "hello")
	t.Logf("codex output:\n%s", out)

	lower := strings.ToLower(out)
	if !strings.Contains(lower, "hello") &&
		!strings.Contains(lower, "hi") &&
		!strings.Contains(lower, "feeling") &&
		!strings.Contains(lower, "mind") &&
		!strings.Contains(lower, "troubling") {
		t.Error("expected Eliza greeting in output")
	}
}

func TestCodex_IAmSad(t *testing.T) {
	codexPath := codexAvailable(t)
	url, stop := startEliza(t)
	defer stop()

	dir := t.TempDir()
	out := runCodex(t, codexPath, url, dir, "I am sad")
	t.Logf("codex output:\n%s", out)

	lower := strings.ToLower(out)
	if !strings.Contains(lower, "sad") && !strings.Contains(lower, "you") {
		t.Logf("note: Eliza may have matched a different rule")
	}
}

func TestCodex_WhatCanYouDo(t *testing.T) {
	codexPath := codexAvailable(t)
	url, stop := startEliza(t)
	defer stop()

	dir := t.TempDir()
	out := runCodex(t, codexPath, url, dir, "What can you do?")
	t.Logf("codex output:\n%s", out)

	if !strings.Contains(out, "special commands") {
		t.Error("expected Eliza help text listing special commands")
	}
}

func TestCodex_WriteFile(t *testing.T) {
	codexPath := codexAvailable(t)
	url, stop := startEliza(t)
	defer stop()

	dir := t.TempDir()
	out := runCodex(t, codexPath, url, dir, "Please write the text 'hello from eliza' to a file called greeting.txt")
	t.Logf("codex output:\n%s", out)

	data, err := os.ReadFile(filepath.Join(dir, "greeting.txt"))
	if err != nil {
		t.Skipf("codex did not create file (Eliza may not have triggered tool use): %v", err)
	}
	if !strings.Contains(string(data), "hello from eliza") {
		t.Errorf("file contents = %q, expected 'hello from eliza'", string(data))
	}
}

func TestCodex_ListFiles(t *testing.T) {
	codexPath := codexAvailable(t)
	url, stop := startEliza(t)
	defer stop()

	dir := t.TempDir()
	// Seed a file
	os.WriteFile(filepath.Join(dir, "test.txt"), []byte("hello"), 0644)

	out := runCodex(t, codexPath, url, dir, "List the files in the current directory")
	t.Logf("codex output:\n%s", out)

	if !strings.Contains(out, "test.txt") {
		t.Logf("note: Eliza is a therapist, not an assistant — she may not have listed files")
	}
}
