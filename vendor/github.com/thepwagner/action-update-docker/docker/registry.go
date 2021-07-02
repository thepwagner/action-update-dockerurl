package docker

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/cli/cli/manifest/types"
	"github.com/docker/cli/cli/registry/client"
	"github.com/docker/distribution/reference"
	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update/version"
)

type TagLister interface {
	// Tags returns potential version tags given a updater.Dependency path
	Tags(ctx context.Context, path string) ([]string, error)
}

type ImagePinner interface {
	// Pin normalizes Docker image name to sha256 pinned image.
	Pin(ctx context.Context, image string) (string, error)
	Unpin(ctx context.Context, image, hash string) (string, error)
}

type RemoteRegistries struct {
	rt http.RoundTripper
}

func NewRemoteRegistries() *RemoteRegistries {
	return &RemoteRegistries{
		rt: http.DefaultTransport,
	}
}

func (r *RemoteRegistries) Tags(ctx context.Context, image string) ([]string, error) {
	// Normalize image name:
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return nil, fmt.Errorf("invalid image name: %w", err)
	}
	logrus.WithField("image", normalized.String()).Debug("listing image tags")

	cli, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}
	if err := cli.Initialize(flags.NewClientOptions()); err != nil {
		return nil, fmt.Errorf("initializing cli: %w", err)
	}

	tags, err := cli.RegistryClient(false).GetTags(ctx, normalized)
	if err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}
	return tags, nil
}

func (r *RemoteRegistries) Pin(ctx context.Context, image string) (string, error) {
	// Normalize image name:
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return "", fmt.Errorf("invalid image name: %w", err)
	}
	logrus.WithField("image", normalized.String()).Debug("listing image tags")

	cli, err := command.NewDockerCli()
	if err != nil {
		return "", err
	}
	if err := cli.Initialize(flags.NewClientOptions()); err != nil {
		return "", fmt.Errorf("initializing cli: %w", err)
	}

	registryClient := cli.RegistryClient(false)
	mf, err := r.getManifest(ctx, registryClient, normalized)
	if err != nil {
		return "", fmt.Errorf("getting manifest: %w", err)
	}
	return mf.Descriptor.Digest.String(), nil
}

func (r *RemoteRegistries) Unpin(ctx context.Context, image, hash string) (string, error) {
	normalized, err := reference.ParseNormalizedNamed(fmt.Sprintf("%s@%s", image, hash))
	if err != nil {
		return "", fmt.Errorf("invalid image name: %w", err)
	}

	cli, err := command.NewDockerCli()
	if err != nil {
		return "", err
	}
	if err := cli.Initialize(flags.NewClientOptions()); err != nil {
		return "", fmt.Errorf("initializing cli: %w", err)
	}

	registryClient := cli.RegistryClient(false)
	tags, err := registryClient.GetTags(ctx, normalized)
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "", fmt.Errorf("tag not found")
	}

	// Filter semver tags, work backwards (assuming the pinned sha is a near-latest version)
	semverTags := make([]string, 0)
	for _, tag := range tags {
		if version.Semverish(tag) == "" {
			continue
		}
		semverTags = append(semverTags, tag)
	}
	semverTags = version.SemverSort(semverTags)

	logrus.WithFields(logrus.Fields{
		"image": normalized.String(),
		"hash":  hash,
		"tags":  len(semverTags),
	}).Info("listing tags to identify SHA")

	for _, tag := range semverTags {
		normalizedTag, err := reference.ParseNormalizedNamed(fmt.Sprintf("%s:%s", normalized.Name(), tag))
		if err != nil {
			continue
		}
		mf, err := r.getManifest(ctx, registryClient, normalizedTag)
		if err != nil {
			continue
		}
		digest := mf.Descriptor.Digest.String()
		logrus.WithFields(logrus.Fields{
			"tag":    tag,
			"digest": digest,
		}).Debug("fetched image details")

		if digest == hash {
			logrus.WithFields(logrus.Fields{
				"digest": digest,
				"tag":    tag,
			}).Info("resolved pinned image to tag")
			return tag, nil
		}
	}
	return "", fmt.Errorf("manifest not found")
}

func (r *RemoteRegistries) getManifest(ctx context.Context, registryClient client.RegistryClient, normalized reference.Named) (*types.ImageManifest, error) {
	// Assume this image is available for one platform:
	mf, err := registryClient.GetManifest(ctx, normalized)
	if err == nil {
		return &mf, nil
	}
	if !strings.Contains(err.Error(), "is a manifest list") {
		return nil, fmt.Errorf("getting manifest: %w", err)
	}

	// Multi-platform images have a list of manifests, select the "right" one:
	manifestList, err := registryClient.GetManifestList(ctx, normalized)
	if err != nil {
		return nil, fmt.Errorf("fetching manifest list: %w", err)
	}
	for _, mf := range manifestList {
		pl := mf.Descriptor.Platform
		if pl.Architecture != "amd64" && pl.OS != "linux" {
			continue
		}
		return &mf, nil
	}
	return nil, fmt.Errorf("could not resolve %q", normalized.String())
}
