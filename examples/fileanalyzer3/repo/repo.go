package repo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CloneRepo clones a git repo to temp dir
func CloneRepo(repo string) (string, error) {
	tmp, err := os.MkdirTemp("", "pocketflow_repo_")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("git", "clone", "--depth", "1", repo, tmp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmp)
		return "", err
	}
	return tmp, nil
}

// CollectFiles recursively collects files matching include/exclude
func CollectFiles(root string, includes, excludes []string, maxSize int64) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			for _, ex := range excludes {
				if strings.HasPrefix(path, filepath.Join(root, ex)) || strings.Contains(path, ex) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if info.Size() > maxSize {
			return nil
		}
		base := filepath.Base(path)
		for _, inc := range includes {
			if ok, _ := filepath.Match(inc, base); ok || base == inc {
				files = append(files, path)
				break
			}
		}
		return nil
	})
	return files, err
}
