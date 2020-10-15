package docker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/distribution/reference"
	"github.com/sirupsen/logrus"
)

type TagLister interface {
	Tags(ctx context.Context, path string) ([]string, error)
}

type RemoteTagLister struct {
	rt http.RoundTripper
}

func NewRemoteTagLister() *RemoteTagLister {
	return &RemoteTagLister{
		rt: http.DefaultTransport,
	}
}

func (r *RemoteTagLister) Tags(ctx context.Context, image string) ([]string, error) {
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
