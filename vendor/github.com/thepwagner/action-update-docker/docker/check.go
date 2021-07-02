package docker

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update/updater"
	"github.com/thepwagner/action-update/version"
	"golang.org/x/mod/semver"
)

func (u *Updater) Check(ctx context.Context, dependency updater.Dependency, filter func(string) bool) (*updater.Update, error) {
	previous := version.Semverish(dependency.Version)
	if previous == "" {
		if !sha256Ish(dependency.Version) {
			logrus.WithFields(logrus.Fields{"path": dependency.Path, "version": dependency.Version}).Debug("ignoring non-semver dependency")
			return nil, nil
		}
		tag, err := u.pinner.Unpin(ctx, dependency.Path, dependency.Version)
		if err != nil {
			return nil, fmt.Errorf("unpinning %q: %w", dependency.Path, err)
		}
		previous = tag
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

		mapped := version.Semverish(t)
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

	versions = version.SemverSort(versions)
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


var sha256VersionRE = regexp.MustCompile("sha256:[a-f0-9]{64}")

func sha256Ish(s string) bool {
	return sha256VersionRE.MatchString(s)
}
