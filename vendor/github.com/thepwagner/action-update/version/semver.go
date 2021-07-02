package version

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

// Semverish coerces semver-ish strings to semver, or returns the empty string
func Semverish(s string) string {
	if v := isValidWithV(s); v != "" {
		return v
	}

	dots := strings.Split(s, ".")
	if len(dots) > 3 {
		normal := strings.Join(dots[:3], ".") + "-" + strings.Join(dots[3:], ".")
		if v := isValidWithV(normal); v != "" {
			return v
		}
	}

	return ""
}

func isValidWithV(s string) string {
	if semver.IsValid(s) {
		return s
	}

	if vt := fmt.Sprintf("v%s", s); semver.IsValid(vt) {
		return vt
	}
	return ""
}

// SemverSort is an descending sort of semver-ish version strings. The latest version is at index 0.
func SemverSort(versions []string) []string {
	sort.Slice(versions, func(i, j int) bool {
		// Prefer strict semver ordering:
		if c := semver.Compare(Semverish(versions[i]), Semverish(versions[j])); c > 0 {
			return true
		} else if c < 0 {
			return false
		}
		// Failing that, prefer the most specific version:
		return strings.Count(versions[i], ".") > strings.Count(versions[j], ".")
	})
	return versions
}
