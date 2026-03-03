package main

import (
	"testing"
)

func TestDefaultLimit(t *testing.T) {
	if defaultLimit <= 0 {
		t.Errorf("expected positive default limit, got %d", defaultLimit)
	}
}
