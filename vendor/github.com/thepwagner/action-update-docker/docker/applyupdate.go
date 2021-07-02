package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/command"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/thepwagner/action-update/updater"
	"github.com/thepwagner/action-update/version"
	"golang.org/x/mod/semver"
)

func (u *Updater) ApplyUpdate(ctx context.Context, update updater.Update) error {
	var nextVersion string
	if u.pinImageSha {
		pinned, err := u.pinner.Pin(ctx, fmt.Sprintf("%s:%s", update.Path, update.Next))
		if err != nil {
			return fmt.Errorf("pinning image: %w", err)
		}
		nextVersion = pinned
	} else {
		nextVersion = update.Next
	}

	return WalkDockerfiles(u.root, u.pathFilter, func(path string, parsed *parser.Result) error {
		vars := NewInterpolation(parsed)

		// prepare a strings.NewReplacer to find/replace strings based on the parsed AST:
		var oldnew []string
		var seenFrom bool
		for _, instruction := range parsed.AST.Children {
			switch instruction.Value {
			case command.From:
				seenFrom = true
				dep := parseDependency(vars, instruction.Next.Value)
				if dep == nil || dep.Path != update.Path {
					continue
				}
				// Ignore FROM statements with a variable:
				if !strings.Contains(instruction.Original, "$") {
					re := regexp.MustCompile(fmt.Sprintf(`%s[:@][^\s]*`, regexp.QuoteMeta(update.Path)))
					var replacement string
					if u.pinImageSha {
						replacement = fmt.Sprintf("%s@%s", update.Path, nextVersion)
					} else {
						replacement = fmt.Sprintf("%s:%s", update.Path, nextVersion)
					}
					newInstruction := re.ReplaceAllString(instruction.Original, replacement)

					oldVersion := fmt.Sprintf("%s:%s", update.Path, update.Previous)
					newVersion := fmt.Sprintf("%s:%s", update.Path, update.Next)
					var commentFound bool
					for _, comment := range instruction.PrevComment {
						if strings.Contains(comment, oldVersion) {
							comment = fmt.Sprintf("# %s", comment)
							oldnew = append(oldnew, comment, re.ReplaceAllString(comment, newVersion))
							commentFound = true
						}
					}
					if u.pinImageSha && !commentFound {
						newInstruction = fmt.Sprintf("# %s\n%s", newVersion, newInstruction)
					}
					oldnew = append(oldnew, instruction.Original, newInstruction)
				}
			case command.Arg:
				if seenFrom {
					continue
				}

				varSplit := strings.SplitN(instruction.Next.Value, "=", 2)
				if len(varSplit) != 2 {
					continue
				}

				if varSplit[1] == update.Previous {
					// Variable is exact version, direct replace
					oldnew = append(oldnew, instruction.Original, strings.ReplaceAll(instruction.Original, update.Previous, nextVersion))
					continue
				}

				suffix := semver.Prerelease(version.Semverish(update.Previous))
				if suffix != "" {
					c := semver.Canonical(version.Semverish(varSplit[1]))
					noSuffix := update.Previous[:len(update.Previous)-len(suffix)]

					if semver.Compare(c, semver.Canonical(version.Semverish(noSuffix))) == 0 {
						nextSuffix := semver.Prerelease(version.Semverish(update.Next))
						nextNoSuffix := update.Next[:len(update.Next)-len(nextSuffix)]
						oldnew = append(oldnew, instruction.Original, strings.ReplaceAll(instruction.Original, noSuffix, nextNoSuffix))
						continue
					}
				}
			}
		}

		if len(oldnew) == 0 {
			return nil
		}

		// Read file into memory:
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Rewrite contents through replacer:
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0640)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := strings.NewReplacer(oldnew...).WriteString(f, string(b)); err != nil {
			return err
		}
		return nil
	})
}
