package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
)

// User represents the authenticated GitHub user profile.
type User struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type keyPayload struct {
	Title string `json:"title"`
	Key   string `json:"key"`
}

// GetCurrentUser returns the profile of the currently authenticated GitHub user.
func GetCurrentUser() (*User, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GitHub API client: %w", err)
	}

	var user User
	if err := client.Get("user", &user); err != nil {
		return nil, fmt.Errorf("failed fetching authenticated user from GitHub API: %w", err)
	}

	return &user, nil
}

// RegisterAuthKey registers an SSH public key for Git transport authentication (push/pull).
func RegisterAuthKey(title, pubKey string) (bool, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return false, fmt.Errorf("failed to initialize GitHub API client: %w", err)
	}

	payload, err := json.Marshal(keyPayload{
		Title: title,
		Key:   pubKey,
	})
	if err != nil {
		return false, err
	}

	err = client.Post("user/keys", bytes.NewReader(payload), nil)
	if err != nil {
		var httpErr *api.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == 422 {
			// Key is already present on the account
			if strings.Contains(strings.ToLower(httpErr.Message), "already") ||
				strings.Contains(strings.ToLower(httpErr.Error()), "already") {
				return false, nil
			}
		}
		return false, fmt.Errorf("failed registering authentication key: %w", err)
	}

	return true, nil
}

// RegisterSigningKey registers an SSH public key for commit/tag cryptographic verification.
func RegisterSigningKey(title, pubKey string) (bool, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return false, fmt.Errorf("failed to initialize GitHub API client: %w", err)
	}

	payload, err := json.Marshal(keyPayload{
		Title: title,
		Key:   pubKey,
	})
	if err != nil {
		return false, err
	}

	err = client.Post("user/ssh_signing_keys", bytes.NewReader(payload), nil)
	if err != nil {
		var httpErr *api.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == 422 {
			// Key is already present on the account
			if strings.Contains(strings.ToLower(httpErr.Message), "already") ||
				strings.Contains(strings.ToLower(httpErr.Error()), "already") {
				return false, nil
			}
		}
		return false, fmt.Errorf("failed registering signing key: %w", err)
	}

	return true, nil
}
