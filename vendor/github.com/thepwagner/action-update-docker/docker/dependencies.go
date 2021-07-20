package docker

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/command"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/thepwagner/action-update/updater"
	"github.com/thepwagner/action-update/version"
)

func (u *Updater) Dependencies(_ context.Context) ([]updater.Dependency, error) {
	var deps []updater.Dependency
	err := WalkDockerfiles(u.root, u.pathFilter, func(path string, parsed *parser.Result) error {
		fileDeps, err := extractImages(parsed)
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

func extractImages(parsed *parser.Result) ([]updater.Dependency, error) {
	vars := NewInterpolation(parsed)

	var deps []updater.Dependency
	for _, instruction := range parsed.AST.Children {
		// Ignore everything but FROM instructions
		if !strings.EqualFold(instruction.Value, command.From) {
			continue
		}

		// Parse the image name:
		image := instruction.Next.Value
		dep := parseDependency(vars, image)
		if dep != nil {
			deps = append(deps, *dep)
		}
	}
	return deps, nil
}

var sha256RE = regexp.MustCompile("[a-f0-9]{64}")

func parseDependency(vars *Interpolation, image string) *updater.Dependency {
	imageSplit := strings.SplitN(image, ":", 2)
	if len(imageSplit) == 1 {
		// No tag provided, default to ":latest"
		return &updater.Dependency{Path: image, Version: "latest"}
	}

	if strings.Contains(imageSplit[1], "$") {
		// Version contains a variable, attempt interpolation:
		vers := vars.Interpolate(imageSplit[1])
		if !strings.Contains(vers, "$") {
			return &updater.Dependency{Path: imageSplit[0], Version: vers}
		}
	} else if version.Semverish(imageSplit[1]) != "" {
		// Image tag is valid semver:
		return &updater.Dependency{Path: imageSplit[0], Version: imageSplit[1]}
	} else if strings.HasSuffix(imageSplit[0], "@sha256") && sha256RE.MatchString(imageSplit[1]) {
		return &updater.Dependency{Path: imageSplit[0][:len(imageSplit[0])-7], Version: fmt.Sprintf("sha256:%s", imageSplit[1])}
	}
	return nil
}
