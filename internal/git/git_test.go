package git

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected Version
		hasError bool
	}{
		{
			input:    "git version 2.54.0.windows.1",
			expected: Version{Major: 2, Minor: 54, Patch: 0},
			hasError: false,
		},
		{
			input:    "git version 2.36.0",
			expected: Version{Major: 2, Minor: 36, Patch: 0},
			hasError: false,
		},
		{
			input:    "git version 2.35.9",
			expected: Version{Major: 2, Minor: 35, Patch: 9},
			hasError: false,
		},
		{
			input:    "invalid",
			expected: Version{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		got, err := ParseVersion(tt.input)
		if tt.hasError && err == nil {
			t.Errorf("ParseVersion(%q) expected error, got nil", tt.input)
		}
		if !tt.hasError && err != nil {
			t.Errorf("ParseVersion(%q) unexpected error: %v", tt.input, err)
		}
		if got != tt.expected {
			t.Errorf("ParseVersion(%q) = %v, expected %v", tt.input, got, tt.expected)
		}
	}
}

func TestVersionGreaterOrEqual(t *testing.T) {
	min := Version{Major: 2, Minor: 36, Patch: 0}

	v1 := Version{Major: 2, Minor: 54, Patch: 0}
	if !v1.GreaterOrEqual(min) {
		t.Error("expected 2.54.0 >= 2.36.0")
	}

	v2 := Version{Major: 2, Minor: 36, Patch: 0}
	if !v2.GreaterOrEqual(min) {
		t.Error("expected 2.36.0 >= 2.36.0")
	}

	v3 := Version{Major: 2, Minor: 35, Patch: 9}
	if v3.GreaterOrEqual(min) {
		t.Error("expected 2.35.9 < 2.36.0")
	}

	v4 := Version{Major: 1, Minor: 99, Patch: 99}
	if v4.GreaterOrEqual(min) {
		t.Error("expected 1.99.99 < 2.36.0")
	}
}
