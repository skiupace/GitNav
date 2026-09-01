package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitnav-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	_, err = NewRepo(tempDir)
	if err == nil {
		t.Errorf("Expected error for non-git repository, got nil")
	}

	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatal(err)
	}

	repo, err := NewRepo(tempDir)
	if err != nil {
		t.Errorf("Expected success for git repository, got error: %v", err)
	}
	if repo.Root != tempDir {
		t.Errorf("Expected root %s, got %s", tempDir, repo.Root)
	}
}
