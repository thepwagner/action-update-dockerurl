package docker

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update/updater"
	"golang.org/x/mod/semver"
)

func (u *Updater) Check(ctx context.Context, dependency updater.Dependency, filter func(string) bool) (*updater.Update, error) {
	previous := semverIsh(dependency.Version)
	if previous == "" {
		logrus.WithFields(logrus.Fields{"path": dependency.Path, "version": dependency.Version}).Debug("ignoring non-semver dependency")
		return nil, nil
	}
	suffix := semver.Prerelease(previous)

	tags, err := u.tags.Tags(ctx, dependency.Path)
	if err != nil {
		return nil, fmt.Errorf("querying tags: %w", err)
	}

	versions := make([]string, 0, len(tags))
	versionMap := map[string]string{}
	for _, t := range tags {
		// Skip datestamped tags as valid-not-valid semver
		if len(t) == 8 && strings.HasPrefix(t, "20") {
			continue
		}

		mapped := semverIsh(t)
		if mapped == "" {
			continue
		}
		if semver.Prerelease(mapped) != suffix {
			continue
		}
		if filter != nil && !filter(mapped) {
			continue
		}

		versions = append(versions, mapped)
		versionMap[mapped] = t
	}
	if len(versions) == 0 {
		return nil, nil
	}

	sort.Slice(versions, func(i, j int) bool {
		// Prefer strict semver ordering:
		if c := semver.Compare(versions[i], versions[j]); c > 0 {
			return true
		} else if c < 0 {
			return false
		}
		// Failing that, prefer the most specific version:
		return strings.Count(versions[i], ".") > strings.Count(versions[j], ".")
	})
	latest := versions[0]
	if semver.Compare(previous, latest) >= 0 {
		return nil, nil
	}

	return &updater.Update{
		Path:     dependency.Path,
		Previous: dependency.Version,
		Next:     versionMap[latest],
	}, nil
}

func semverIsh(s string) string {
	if semver.IsValid(s) {
		return s
	}

	if vt := fmt.Sprintf("v%s", s); semver.IsValid(vt) {
		return vt
	}
	return ""
}
