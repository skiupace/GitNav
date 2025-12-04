package git

import (
	"os"
	"path/filepath"
	"sort"

	ignore "github.com/sabhiram/go-gitignore"
)

type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*Node
}

// ScanRepo returns a tree of files/folders under repo.Root
func (r *Repo) ScanRepo() (*Node, error) {
	gitignorePath := filepath.Join(r.Root, ".gitignore")
	var ign *ignore.GitIgnore
	if _, err := os.Stat(gitignorePath); err == nil {
		ign, _ = ignore.CompileIgnoreFile(gitignorePath)
	}

	return scanDir(r.Root, ign, r.Root)
}

func scanDir(path string, ign *ignore.GitIgnore, repoRoot string) (*Node, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Skip the .git folder
	if info.IsDir() && info.Name() == ".git" {
		return nil, nil
	}

	relPath, _ := filepath.Rel(repoRoot, path)

	// Skip if ignored
	if ign != nil && ign.MatchesPath(relPath) {
		return nil, nil
	}

	node := &Node{
		Name:  info.Name(),
		Path:  path,
		IsDir: info.IsDir(),
	}

	if !info.IsDir() {
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		childNode, err := scanDir(childPath, ign, repoRoot)
		if err != nil || childNode == nil {
			continue
		}
		node.Children = append(node.Children, childNode)
	}

	// If directory has no children after ignoring, skip it entirely
	if len(node.Children) == 0 {
		return nil, nil
	}

	// ---- SORT: folders first, then files ----
	sort.Slice(node.Children, func(i, j int) bool {
		a, b := node.Children[i], node.Children[j]

		if a.IsDir != b.IsDir {
			return a.IsDir // dirs first
		}

		return a.Name < b.Name // alphabetical
	})

	return node, nil
}
