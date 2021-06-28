module github.com/thepwagner/action-update-dockerurl

go 1.15

require (
	github.com/google/go-github/v36 v36.0.0
	github.com/moby/buildkit v0.8.3
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/thepwagner/action-update v0.0.40
	github.com/thepwagner/action-update-docker v0.0.8
	golang.org/x/mod v0.4.2
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.0
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
)
