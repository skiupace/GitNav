package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	gogit "github.com/go-git/go-git/v5"
)

type Repo struct {
	Root string
}

func NewRepo(path string) (*Repo, error) {
	// .git may be a directory (normal repo) or a file (worktree/submodule)
	if _, err := os.Stat(filepath.Join(path, ".git")); err != nil {
		return nil, ErrNotAGitRepo
	}

	return &Repo{Root: path}, nil
}

var ErrNotAGitRepo = fmt.Errorf("not a git repository")

// Info holds best-effort repository metadata.
type Info struct {
	Branch  string
	SHA     string
	Subject string
	Author  string
	When    time.Time
	Dirty   int
}

// Info returns HEAD and worktree status. Errors are swallowed.
// fields simply stay zero-valued when unavailable.
func (r *Repo) Info() Info {
	var in Info

	repo, err := gogit.PlainOpen(r.Root)
	if err != nil {
		return in
	}

	ref, err := repo.Head()
	if err != nil {
		in.Branch = "no commits yet"
	} else {
		in.Branch = ref.Name().Short()
		in.SHA = ref.Hash().String()[:7]
		if c, err := repo.CommitObject(ref.Hash()); err == nil {
			in.Subject = strings.SplitN(strings.TrimSpace(c.Message), "\n", 2)[0]
			in.Author = c.Author.Name
			in.When = c.Author.When
		}
	}

	if w, err := repo.Worktree(); err == nil {
		if st, err := w.Status(); err == nil {
			for _, f := range st {
				if f.Staging != gogit.Unmodified || f.Worktree != gogit.Unmodified {
					in.Dirty++
				}
			}
		}
	}

	return in
}
