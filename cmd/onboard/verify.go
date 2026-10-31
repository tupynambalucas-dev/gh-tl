package onboard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tupynambalucas-dev/gh-tl/internal/git"
	"github.com/tupynambalucas-dev/gh-tl/internal/github"
	"github.com/tupynambalucas-dev/gh-tl/internal/ssh"
	"github.com/tupynambalucas-dev/gh-tl/internal/ui"
)

// NewVerifyCmd creates and configures the 'verify' cobra command.
func NewVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Validate workstation compliance across the 7 pre-flight readiness gates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVerify()
		},
	}
}

// RunVerify executes the 7 compliance checks without altering local configuration.
func RunVerify() error {
	ui.PrintHeader("tupynambalucas-dev Workstation Readiness Gate Verification")

	failedGates := 0

	// Gate 1: Git Version Check (>= 2.36.0)
	fmt.Print("[Gate 1/7] Inspecting Git version... ")
	ver, raw, err := git.GetGitVersion()
	minVer := git.Version{Major: 2, Minor: 36, Patch: 0}
	if err != nil {
		ui.PrintFail("Git not found in PATH")
		ui.PrintRemediation("Install Git >= 2.36.0 via 'winget install Git.Git' or 'brew install git'")
		failedGates++
	} else if ver.GreaterOrEqual(minVer) {
		ui.PrintPass(fmt.Sprintf("%s (>= 2.36.0)", raw))
	} else {
		ui.PrintFail(fmt.Sprintf("Git %s detected. Minimum required is 2.36.0", ver.String()))
		ui.PrintRemediation("Upgrade Git via 'winget upgrade Git.Git' or 'brew upgrade git'")
		failedGates++
	}

	// Gate 2: SSH Signing Key Existence
	fmt.Print("[Gate 2/7] Inspecting Ed25519 signing key... ")
	home, _ := os.UserHomeDir()
	defaultPubKey := filepath.Join(home, ".ssh", "github_tupynambalucas-dev_ed25519.pub")
	legacyPubKey := filepath.Join(home, ".ssh", "github_tupynambalucas_ed25519.pub")
	signingKeyConfig, _ := git.GetConfig("user.signingkey")

	if signingKeyConfig != "" && ssh.KeyExists(signingKeyConfig) {
		ui.PrintPass(fmt.Sprintf("Active signing key detected: %s", signingKeyConfig))
	} else if ssh.KeyExists(defaultPubKey) {
		ui.PrintPass(fmt.Sprintf("Key found at %s", defaultPubKey))
	} else if ssh.KeyExists(legacyPubKey) {
		ui.PrintPass(fmt.Sprintf("Key found at %s", legacyPubKey))
	} else {
		ui.PrintFail("Ed25519 signing key not found")
		ui.PrintRemediation("Run 'gh tl onboard' to generate and provision your corporate keypair")
		failedGates++
	}

	// Gate 3: Commit Cryptographic Signing Configuration
	fmt.Print("[Gate 3/7] Inspecting cryptographic commit signing config... ")
	gpgFormat, _ := git.GetConfig("gpg.format")
	commitGpgSign, _ := git.GetConfig("commit.gpgsign")
	if strings.EqualFold(gpgFormat, "ssh") && strings.EqualFold(commitGpgSign, "true") {
		ui.PrintPass("gpg.format=ssh, commit.gpgsign=true")
	} else {
		ui.PrintFail(fmt.Sprintf("Signing not configured globally (gpg.format=%s, commit.gpgsign=%s)", gpgFormat, commitGpgSign))
		ui.PrintRemediation("Execute: git config --global gpg.format ssh && git config --global commit.gpgsign true")
		failedGates++
	}

	// Gate 4: Mani Multi-Repository CLI Orchestrator
	fmt.Print("[Gate 4/7] Inspecting Mani multi-repository orchestrator... ")
	maniOut, err := git.CheckCommand("mani", "--version")
	if err == nil {
		firstLine := strings.Split(maniOut, "\n")[0]
		ui.PrintPass(firstLine)
	} else {
		ui.PrintFail("Mani CLI orchestrator not found in PATH")
		ui.PrintRemediation("Install Mani via 'scoop install mani' or download from https://github.com/aliesbelik/mani/releases")
		failedGates++
	}

	// Gate 5: Node.js & pnpm Environment
	fmt.Print("[Gate 5/7] Inspecting Node.js and pnpm runtime... ")
	nodeOut, errNode := git.CheckCommand("node", "--version")
	pnpmOut, errPnpm := git.CheckCommand("pnpm", "--version")
	if errNode == nil && errPnpm == nil {
		ui.PrintPass(fmt.Sprintf("Node %s, pnpm %s", strings.TrimSpace(nodeOut), strings.TrimSpace(pnpmOut)))
	} else {
		ui.PrintFail(fmt.Sprintf("Runtime missing (node: %v, pnpm: %v)", errNode == nil, errPnpm == nil))
		ui.PrintRemediation("Install Node.js >= 22 (LTS) and enable pnpm via 'corepack enable'")
		failedGates++
	}

	// Gate 6: Git SSH Protocol Default
	fmt.Print("[Gate 6/7] Inspecting default Git SSH transport override... ")
	urlInsteadOf, _ := git.GetConfig("url.git@github.com:.insteadOf")
	if strings.EqualFold(urlInsteadOf, "https://github.com/") {
		ui.PrintPass("SSH transport rewrite configured (url.git@github.com:.insteadOf = https://github.com/)")
	} else {
		ui.PrintWarn(fmt.Sprintf("url.git@github.com:.insteadOf is '%s' (recommended: 'https://github.com/')", urlInsteadOf))
		ui.PrintRemediation("Execute: git config --global url.\"git@github.com:\".insteadOf https://github.com/")
	}

	// Gate 7: GitHub CLI Authentication
	fmt.Print("[Gate 7/7] Inspecting GitHub CLI authentication... ")
	user, err := github.GetCurrentUser()
	if err == nil && user.Login != "" {
		ui.PrintPass(fmt.Sprintf("Authenticated as @%s (%s)", user.Login, user.Name))
	} else {
		cmdAuth := exec.Command("gh", "auth", "status")
		outAuth, errAuth := cmdAuth.CombinedOutput()
		if errAuth == nil && strings.Contains(string(outAuth), "Logged in to github.com") {
			ui.PrintPass("GitHub CLI authenticated")
		} else {
			ui.PrintFail("GitHub CLI is not authenticated")
			ui.PrintRemediation("Execute: gh auth login --git-protocol ssh --web")
			failedGates++
		}
	}

	fmt.Println()
	fmt.Println(ui.Gray("----------------------------------------------------------"))

	if failedGates == 0 {
		ui.PrintPass("Workstation is fully compliant with all readiness gates.")
		return nil
	}

	return fmt.Errorf("readiness verification failed: %d gate(s) did not pass", failedGates)
}
