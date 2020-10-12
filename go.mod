module github.com/thepwagner/action-update-dockerurl

go 1.15

require (
	github.com/google/go-github/v32 v32.1.0
	github.com/moby/buildkit v0.7.2
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/thepwagner/action-update v0.0.11
	github.com/thepwagner/action-update-docker v0.0.5
	golang.org/x/mod v0.3.0
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.0
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
)
