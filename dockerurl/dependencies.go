package dockerurl

import (
	"context"
	"fmt"
	"regexp"

	"github.com/moby/buildkit/frontend/dockerfile/command"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/thepwagner/action-update-docker/docker"
	"github.com/thepwagner/action-update/updater"
)

func (u *Updater) Dependencies(_ context.Context) ([]updater.Dependency, error) {
	var deps []updater.Dependency
	err := docker.WalkDockerfiles(u.root, u.pathFilter, func(path string, parsed *parser.Result) error {
		fileDeps, err := extractDockerfile(parsed)
		if err != nil {
			return fmt.Errorf("extracting dependencies: %w", err)
		}
		deps = append(deps, fileDeps...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("collecting dependencies: %w", err)
	}
	return deps, nil
}

var ghRelease = regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+)/releases/download/([^/]+)/`)

func extractDockerfile(parsed *parser.Result) ([]updater.Dependency, error) {
	vars := docker.NewInterpolation(parsed)

	deps := make([]updater.Dependency, 0)
	for _, instruction := range parsed.AST.Children {
		// Ignore everything but RUN instructions
		if instruction.Value != command.Run {
			continue
		}

		// Best-effort interpolate, then extract GitHub release URLs from the resulting commands:
		cmdLine := vars.Interpolate(instruction.Next.Value)
		for _, ghReleaseMatch := range ghRelease.FindAllStringSubmatch(cmdLine, -1) {
			repo := ghReleaseMatch[1]
			name := ghReleaseMatch[2]
			vers := ghReleaseMatch[3]
			deps = append(deps, updater.Dependency{
				Path:    formatGitHubRelease(repo, name),
				Version: vers,
			})
		}
	}
	return deps, nil
}
