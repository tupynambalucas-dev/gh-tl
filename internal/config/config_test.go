package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	windowsPath := `C:\Users\developer\.ssh\id_ed25519.pub`
	norm := NormalizePath(windowsPath)
	if strings.Contains(norm, `\`) {
		t.Errorf("NormalizePath should not contain backslashes, got: %s", norm)
	}
	if norm != "C:/Users/developer/.ssh/id_ed25519.pub" {
		t.Errorf("unexpected normalized path: %s", norm)
	}
}

func TestSetupCorporateGitConfig(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, ".gitconfig-test")
	pubKeyPath := filepath.Join(tempDir, "key.pub")
	allowedSignersPath := filepath.Join(tempDir, "allowed_signers")

	err := SetupCorporateGitConfig("Test Developer", "dev@example.com", pubKeyPath, allowedSignersPath, targetPath)
	if err != nil {
		t.Fatalf("SetupCorporateGitConfig returned error: %v", err)
	}

	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read generated config: %v", err)
	}

	strContent := string(content)
	if !strings.Contains(strContent, "name = Test Developer") {
		t.Error("config missing user name")
	}
	if !strings.Contains(strContent, "email = dev@example.com") {
		t.Error("config missing user email")
	}
	if !strings.Contains(strContent, "gpgsign = true") {
		t.Error("config missing gpgsign = true")
	}
	if !strings.Contains(strContent, "format = ssh") {
		t.Error("config missing format = ssh")
	}
	if strings.Contains(strContent, `\`) {
		t.Error("config contains un-normalized backslashes")
	}
}
