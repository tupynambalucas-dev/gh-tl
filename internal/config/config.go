package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tupynambalucas-dev/gh-tl/internal/git"
)

// NormalizePath converts backslashes to forward slashes for Git config compatibility.
func NormalizePath(p string) string {
	return strings.ReplaceAll(filepath.Clean(p), "\\", "/")
}

// ApplyGlobalGitDefaults configures standard baseline Git configurations.
func ApplyGlobalGitDefaults() error {
	defaults := []struct {
		key   string
		value string
	}{
		{"init.defaultBranch", "main"},
		{"url.git@github.com:.insteadOf", "https://github.com/"},
		{"gpg.format", "ssh"},
		{"commit.gpgsign", "true"},
		{"tag.gpgsign", "true"},
	}

	for _, d := range defaults {
		if err := git.SetConfigGlobal(d.key, d.value); err != nil {
			return err
		}
	}

	// Platform specific autocrlf
	if runtime.GOOS == "windows" {
		if err := git.SetConfigGlobal("core.autocrlf", "true"); err != nil {
			return err
		}
	} else {
		if err := git.SetConfigGlobal("core.autocrlf", "input"); err != nil {
			return err
		}
	}

	return nil
}

// SetupCorporateGitConfig writes the isolated ~/.gitconfig-tupynambalucas-dev file.
func SetupCorporateGitConfig(userName, userEmail, pubKeyPath, allowedSignersPath, targetPath string) error {
	normPubKey := NormalizePath(pubKeyPath)
	normAllowedSigners := NormalizePath(allowedSignersPath)

	content := fmt.Sprintf(`[user]
	name = %s
	email = %s
	signingkey = %s

[commit]
	gpgsign = true

[tag]
	gpgsign = true

[gpg]
	format = ssh

[gpg "ssh"]
	allowedSignersFile = %s
`, userName, userEmail, normPubKey, normAllowedSigners)

	if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write corporate git config to %s: %w", targetPath, err)
	}

	return nil
}

// SetupGlobalConditionalIncludes binds the corporate git config to tupynambalucas-dev repositories.
func SetupGlobalConditionalIncludes(corpConfigPath string, orgName string) error {
	normCorpPath := NormalizePath(corpConfigPath)

	// Pattern 1: SSH Remote URL (git@github.com:org/**)
	sshKey := fmt.Sprintf("includeIf.hasconfig:remote.*.url:git@github.com:%s/**.path", orgName)
	if err := git.SetConfigGlobal(sshKey, normCorpPath); err != nil {
		return fmt.Errorf("failed setting ssh conditional include: %w", err)
	}

	// Pattern 2: HTTPS Remote URL (https://github.com/org/**)
	httpsKey := fmt.Sprintf("includeIf.hasconfig:remote.*.url:https://github.com/%s/**.path", orgName)
	if err := git.SetConfigGlobal(httpsKey, normCorpPath); err != nil {
		return fmt.Errorf("failed setting https conditional include: %w", err)
	}

	return nil
}
