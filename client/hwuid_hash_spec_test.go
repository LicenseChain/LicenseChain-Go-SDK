package client

import (
	"regexp"
	"testing"
)

func TestGenerateDefaultHWUIDHashSpec(t *testing.T) {
	h1 := GenerateDefaultHWUID()
	h2 := GenerateDefaultHWUID()

	if h1 != h2 {
		t.Fatalf("expected deterministic hwuid hash, got %q and %q", h1, h2)
	}
	if !regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(h1) {
		t.Fatalf("expected lowercase sha256 hex, got %q", h1)
	}
}
