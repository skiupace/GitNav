package git

import (
	"fmt"
	"os"
	"path/filepath"
)

type Repo struct {
	Root string
}

func NewRepo(path string) (*Repo, error) {
	gitPath := filepath.Join(path, ".git")

	info, err := os.Stat(gitPath)
	if err != nil || !info.IsDir() {
		return nil, ErrNotAGitRepo
	}

	return &Repo{Root: path}, nil
}

var ErrNotAGitRepo = fmt.Errorf("not a git repository")
