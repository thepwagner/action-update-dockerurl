module github.com/thepwagner/action-update-dockerurl

go 1.15

require (
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/Microsoft/hcsshim v0.8.18 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210512092938-c05353c2d58c // indirect
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/containerd/continuity v0.1.0 // indirect
	github.com/docker/cli v20.10.7+incompatible // indirect
	github.com/docker/docker v20.10.7+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/google/go-github/v36 v36.0.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/moby/buildkit v0.8.3
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/opencontainers/runc v1.0.0 // indirect
	github.com/prometheus/common v0.29.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/thepwagner/action-update v0.0.42
	github.com/thepwagner/action-update-docker v0.0.9
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/mod v0.4.2
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	google.golang.org/genproto v0.0.0-20210701191553-46259e63a0a9 // indirect
	google.golang.org/grpc v1.40.0-dev // indirect
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.0
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
)
