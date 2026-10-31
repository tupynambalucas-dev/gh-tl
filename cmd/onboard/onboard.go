package onboard

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tupynambalucas-dev/gh-tl/internal/config"
	"github.com/tupynambalucas-dev/gh-tl/internal/git"
	"github.com/tupynambalucas-dev/gh-tl/internal/github"
	"github.com/tupynambalucas-dev/gh-tl/internal/ssh"
	"github.com/tupynambalucas-dev/gh-tl/internal/ui"
)

var (
	flagName       string
	flagEmail      string
	flagKeyPath    string
	flagOrg        string
	flagSkipUpload bool
	flagVerify     bool
)

// NewOnboardCmd creates and configures the 'onboard' cobra command.
func NewOnboardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "onboard",
		Short: "Automated workstation provisioning, SSH key generation, and identity isolation",
		Long: `Configures cryptographic SSH signing, conditional Git identities via
hasconfig:remote.*.url, multi-repo tooling, and registers SSH keys with GitHub API.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagVerify {
				return RunVerify()
			}
			return RunSetup()
		},
	}

	cmd.Flags().StringVar(&flagName, "name", "", "Developer full name for Git commits")
	cmd.Flags().StringVar(&flagEmail, "email", "contato@tupynambalucas.dev", "Corporate email address")
	cmd.Flags().StringVar(&flagKeyPath, "key-path", "", "Target SSH private key path (defaults to ~/.ssh/github_tupynambalucas-dev_ed25519)")
	cmd.Flags().StringVar(&flagOrg, "org", "tupynambalucas-dev", "Target GitHub organization name")
	cmd.Flags().BoolVar(&flagSkipUpload, "skip-upload", false, "Skip registering SSH keys via GitHub API")
	cmd.Flags().BoolVar(&flagVerify, "verify", false, "Execute pre-flight readiness verification only")

	return cmd
}

// RunSetup executes the end-to-end automated workstation bootstrapping workflow.
func RunSetup() error {
	ui.PrintHeader("tupynambalucas-dev Workstation Automated Bootstrapper")

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	// 1. Verify Git version prerequisite
	ui.PrintSubheader("[1/7] Inspecting Git Prerequisites")
	ver, raw, err := git.GetGitVersion()
	if err != nil {
		return fmt.Errorf("git is not installed or not available in PATH: %w", err)
	}
	minVer := git.Version{Major: 2, Minor: 36, Patch: 0}
	if !ver.GreaterOrEqual(minVer) {
		return fmt.Errorf("git %s detected. Version >= 2.36.0 is strictly required for conditional includes and SSH signing", ver.String())
	}
	ui.PrintPass(fmt.Sprintf("%s detected", raw))

	// 2. Resolve Developer Identity
	ui.PrintSubheader("[2/7] Resolving Developer Identity")
	devName := flagName
	devEmail := flagEmail

	ghUser, ghErr := github.GetCurrentUser()
	if ghErr == nil && ghUser != nil {
		ui.PrintInfo(fmt.Sprintf("Authenticated as GitHub user: @%s (%s)", ghUser.Login, ghUser.Name))
		if devName == "" {
			if ghUser.Name != "" {
				devName = ghUser.Name
			} else {
				devName = ghUser.Login
			}
		}
	} else {
		ui.PrintWarn("GitHub CLI is not currently authenticated. Key upload will be skipped if unauthenticated.")
	}

	if devName == "" {
		existingName, _ := git.GetConfig("user.name")
		if existingName != "" {
			devName = existingName
		} else {
			devName = "tupynambalucasdev"
		}
	}
	ui.PrintPass(fmt.Sprintf("Name: %s", devName))
	ui.PrintPass(fmt.Sprintf("Email: %s", devEmail))

	// 3. Provision SSH Ed25519 Keypair
	ui.PrintSubheader("[3/7] Provisioning Corporate Ed25519 SSH Keypair")
	sshDir, err := ssh.EnsureSSHDirectory()
	if err != nil {
		return err
	}

	keyPath := flagKeyPath
	if keyPath == "" {
		keyPath = filepath.Join(sshDir, fmt.Sprintf("github_%s_ed25519", flagOrg))
	}
	pubKeyPath := keyPath + ".pub"

	if !ssh.KeyExists(keyPath) {
		ui.PrintInfo(fmt.Sprintf("Generating Ed25519 keypair at %s...", keyPath))
		if err := ssh.GenerateEd25519Key(keyPath, devEmail); err != nil {
			return err
		}
		ui.PrintPass("SSH keypair generated successfully")
	} else {
		ui.PrintPass(fmt.Sprintf("Existing keypair found at %s", keyPath))
	}

	pubKeyContent, err := ssh.ReadPublicKey(pubKeyPath)
	if err != nil {
		return err
	}

	// 4. Register SSH Key with GitHub API (Authentication + Signing)
	ui.PrintSubheader("[4/7] Synchronizing SSH Key with GitHub API")
	if !flagSkipUpload && ghErr == nil {
		hostname, _ := os.Hostname()
		if hostname == "" {
			hostname = "workstation"
		}

		authTitle := fmt.Sprintf("[%s Auth] %s", flagOrg, hostname)
		createdAuth, errAuth := github.RegisterAuthKey(authTitle, pubKeyContent)
		if errAuth != nil {
			ui.PrintWarn(fmt.Sprintf("Failed registering auth key: %v", errAuth))
		} else if createdAuth {
			ui.PrintPass(fmt.Sprintf("Authentication key registered: '%s'", authTitle))
		} else {
			ui.PrintPass("Authentication key already registered on GitHub profile")
		}

		signTitle := fmt.Sprintf("[%s Sign] %s", flagOrg, hostname)
		createdSign, errSign := github.RegisterSigningKey(signTitle, pubKeyContent)
		if errSign != nil {
			ui.PrintWarn(fmt.Sprintf("Failed registering signing key: %v", errSign))
		} else if createdSign {
			ui.PrintPass(fmt.Sprintf("Signing key registered: '%s'", signTitle))
		} else {
			ui.PrintPass("Signing key already registered on GitHub profile")
		}
	} else {
		ui.PrintInfo("Key upload skipped (flag --skip-upload or unauthenticated)")
	}

	// 5. Setup Isolated Corporate Git Configuration
	ui.PrintSubheader("[5/7] Configuring Isolated Corporate Git Identity")
	corpConfigPath := filepath.Join(home, fmt.Sprintf(".gitconfig-%s", flagOrg))
	allowedSignersPath := filepath.Join(sshDir, "allowed_signers")

	if err := config.SetupCorporateGitConfig(devName, devEmail, pubKeyPath, allowedSignersPath, corpConfigPath); err != nil {
		return err
	}
	ui.PrintPass(fmt.Sprintf("Corporate configuration written to %s", corpConfigPath))

	if err := config.SetupGlobalConditionalIncludes(corpConfigPath, flagOrg); err != nil {
		return err
	}
	ui.PrintPass(fmt.Sprintf("Bound conditional includes via hasconfig:remote.*.url for '%s'", flagOrg))

	// 6. Apply Global Baseline Defaults & Allowed Signers
	ui.PrintSubheader("[6/7] Applying Global Defaults and Trust Store")
	if err := config.ApplyGlobalGitDefaults(); err != nil {
		return err
	}
	ui.PrintPass("Applied standard Git defaults (defaultBranch=main, gpg.format=ssh, commit.gpgsign=true)")

	if err := ssh.UpdateAllowedSigners(allowedSignersPath, devEmail, pubKeyContent); err != nil {
		return err
	}
	ui.PrintPass(fmt.Sprintf("Updated %s with local developer public key", allowedSignersPath))

	// 7. Execute Post-Setup Verification
	ui.PrintSubheader("[7/7] Validating Pre-Flight Verification Gates")
	fmt.Println()
	return RunVerify()
}
