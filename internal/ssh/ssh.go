package ssh

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EnsureSSHDirectory guarantees ~/.ssh exists with secure permissions.
func EnsureSSHDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}

	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create ~/.ssh directory: %w", err)
	}
	return sshDir, nil
}

// KeyExists checks if the given file path exists on disk.
func KeyExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// GenerateEd25519Key executes ssh-keygen without passphrase.
func GenerateEd25519Key(keyPath, comment string) error {
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", comment, "-f", keyPath, "-N", "")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ssh-keygen failed: %s (%w)", strings.TrimSpace(string(out)), err)
	}
	return nil
}

// ReadPublicKey reads the public key content and strips trailing whitespace.
func ReadPublicKey(pubKeyPath string) (string, error) {
	content, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read public key from %s: %w", pubKeyPath, err)
	}
	return strings.TrimSpace(string(content)), nil
}

// UpdateAllowedSigners ensures the user public key is in ~/.ssh/allowed_signers.
func UpdateAllowedSigners(signersPath, email, pubKey string) error {
	dir := filepath.Dir(signersPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory for allowed_signers: %w", err)
	}

	// Extract the key part (e.g. "ssh-ed25519 AAAAC3NzaC1...")
	parts := strings.Fields(pubKey)
	if len(parts) < 2 {
		return fmt.Errorf("invalid public key format: expected at least type and base64 payload")
	}
	keyTypeAndData := fmt.Sprintf("%s %s", parts[0], parts[1])
	newLine := fmt.Sprintf("%s %s", email, keyTypeAndData)

	var existingLines []string
	keyAlreadyPresent := false

	if KeyExists(signersPath) {
		file, err := os.Open(signersPath)
		if err != nil {
			return fmt.Errorf("failed to open allowed_signers: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				existingLines = append(existingLines, line)
				continue
			}

			if strings.Contains(line, keyTypeAndData) {
				keyAlreadyPresent = true
			}
			existingLines = append(existingLines, line)
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed reading allowed_signers: %w", err)
		}
	}

	if keyAlreadyPresent {
		return nil
	}

	existingLines = append(existingLines, newLine)
	content := strings.Join(existingLines, "\n") + "\n"
	if err := os.WriteFile(signersPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write allowed_signers: %w", err)
	}

	return nil
}
