package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version (major, minor, patch).
type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) GreaterOrEqual(other Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch >= other.Patch
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// ParseVersion extracts semver from git --version output.
func ParseVersion(output string) (Version, error) {
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 4 {
		return Version{}, fmt.Errorf("unable to parse git version from: %s", output)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

// GetGitVersion executes git --version and parses the output.
func GetGitVersion() (Version, string, error) {
	cmd := exec.Command("git", "--version")
	out, err := cmd.Output()
	if err != nil {
		return Version{}, "", fmt.Errorf("git is not installed or not in PATH: %w", err)
	}

	raw := strings.TrimSpace(string(out))
	ver, err := ParseVersion(raw)
	if err != nil {
		return Version{}, raw, err
	}
	return ver, raw, nil
}

// GetConfig retrieves a git config value.
func GetConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", "--get", key)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// SetConfigGlobal sets a global git config key/value.
func SetConfigGlobal(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set git config %s=%s: %s (%w)", key, value, stderr.String(), err)
	}
	return nil
}

// CheckCommand checks if a given command executable exists in PATH.
func CheckCommand(name string, args ...string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%s not found in PATH", name)
	}
	if len(args) > 0 {
		cmd := exec.Command(path, args...)
		out, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(out)), nil
		}
	}
	return path, nil
}
