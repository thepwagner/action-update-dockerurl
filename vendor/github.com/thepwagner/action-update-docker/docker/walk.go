package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func WalkDockerfiles(root string, filter func(string) bool, walkFunc func(path string, parsed *parser.Result) error) error {
	err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filter != nil {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if filter(rel) {
				return nil
			}
		}

		if fi.IsDir() || !strings.HasPrefix(filepath.Base(path), "Dockerfile") {
			return nil
		}

		parsed, err := parseDockerfile(path)
		if err != nil {
			return fmt.Errorf("parsing %q: %w", path, err)
		}
		if err := walkFunc(path, parsed); err != nil {
			return fmt.Errorf("walking %q: %w", path, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walking filesystem: %w", err)
	}
	return nil
}

func parseDockerfile(path string) (*parser.Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening dockerfile: %w", err)
	}
	defer f.Close()
	parsed, err := parser.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("parsing dockerfile: %w", err)
	}
	return parsed, nil
}
