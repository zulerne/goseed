package prompt

import "testing"

func TestValidName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"simple", "myapp", true},
		{"with dash", "my-app", true},
		{"with numbers", "app123", true},
		{"dash and numbers", "my-app-2", true},
		{"single char", "a", true},
		{"starts with number", "123app", false},
		{"uppercase", "MyApp", false},
		{"empty", "", false},
		{"underscore", "my_app", false},
		{"with dot", "my.app", false},
		{"starts with dash", "-app", false},
		{"spaces", "my app", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validName.MatchString(tt.input); got != tt.valid {
				t.Errorf("validName.MatchString(%q) = %v, want %v", tt.input, got, tt.valid)
			}
		})
	}
}

func TestInferOwner(t *testing.T) {
	tests := []struct {
		name       string
		modulePath string
		want       string
	}{
		{"github path", "github.com/user/repo", "user"},
		{"short path", "example.com/pkg", "pkg"},
		{"single segment", "mypackage", ""},
		{"empty", "", ""},
		{"deep path", "github.com/org/repo/sub", "org"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferOwner(tt.modulePath); got != tt.want {
				t.Errorf("inferOwner(%q) = %q, want %q", tt.modulePath, got, tt.want)
			}
		})
	}
}

func TestGuessGitHubUser(t *testing.T) {
	// guessGitHubUser should always return a non-empty string
	// (falls back to system username or "user")
	if got := guessGitHubUser(); got == "" {
		t.Error("guessGitHubUser() returned empty string")
	}
}
