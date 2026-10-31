package ssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateAllowedSigners(t *testing.T) {
	tempDir := t.TempDir()
	signersPath := filepath.Join(tempDir, "allowed_signers")
	email := "developer@example.com"
	pubKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI1234567890abcdef comment@host"

	// First run: should create file and add key
	if err := UpdateAllowedSigners(signersPath, email, pubKey); err != nil {
		t.Fatalf("UpdateAllowedSigners failed: %v", err)
	}

	content, err := os.ReadFile(signersPath)
	if err != nil {
		t.Fatalf("failed to read signers file: %v", err)
	}

	expectedPrefix := "developer@example.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI1234567890abcdef"
	if !strings.Contains(string(content), expectedPrefix) {
		t.Errorf("signers file does not contain expected entry: %s", string(content))
	}

	// Second run: idempotency check (must not duplicate)
	if err := UpdateAllowedSigners(signersPath, email, pubKey); err != nil {
		t.Fatalf("second UpdateAllowedSigners failed: %v", err)
	}

	contentSecond, _ := os.ReadFile(signersPath)
	occurrences := strings.Count(string(contentSecond), "AAAAC3NzaC1lZDI1NTE5AAAAI1234567890abcdef")
	if occurrences != 1 {
		t.Errorf("expected exactly 1 occurrence of key, found %d", occurrences)
	}
}
