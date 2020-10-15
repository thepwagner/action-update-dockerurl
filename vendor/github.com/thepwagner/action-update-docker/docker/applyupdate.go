package docker

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/command"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/thepwagner/action-update/updater"
	"golang.org/x/mod/semver"
)

func (u *Updater) ApplyUpdate(_ context.Context, update updater.Update) error {
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
				if !strings.Contains(instruction.Original, "$") {
					oldnew = append(oldnew, instruction.Original, strings.ReplaceAll(instruction.Original, update.Previous, update.Next))
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
					oldnew = append(oldnew, instruction.Original, strings.ReplaceAll(instruction.Original, update.Previous, update.Next))
					continue
				}

				suffix := semver.Prerelease(semverIsh(update.Previous))
				if suffix != "" {
					c := semver.Canonical(semverIsh(varSplit[1]))
					noSuffix := update.Previous[:len(update.Previous)-len(suffix)]

					if semver.Compare(c, semver.Canonical(semverIsh(noSuffix))) == 0 {
						nextSuffix := semver.Prerelease(semverIsh(update.Next))
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
		f, err := os.OpenFile(path, os.O_RDWR, 0640)
		if err != nil {
			return err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		// Rewrite contents through replacer:
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return err
		}
		if _, err := strings.NewReplacer(oldnew...).WriteString(f, string(b)); err != nil {
			return err
		}
		return nil
	})
}
