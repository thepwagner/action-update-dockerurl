package docker

import (
	"github.com/thepwagner/action-update/updater"
)

type Updater struct {
	root        string
	pathFilter  func(string) bool
	pinImageSha bool

	tags   TagLister
	pinner ImagePinner
}

var _ updater.Updater = (*Updater)(nil)

func NewUpdater(root string, opts ...UpdaterOpt) *Updater {
	reg := NewRemoteRegistries()
	u := &Updater{
		root:   root,
		tags:   reg,
		pinner: reg,
	}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func (u *Updater) Name() string { return "docker" }

type UpdaterOpt func(*Updater)

func WithTagsLister(tags TagLister) UpdaterOpt {
	return func(u *Updater) {
		u.tags = tags
	}
}

func WithImagePinner(pinner ImagePinner) UpdaterOpt {
	return func(u *Updater) {
		u.pinner = pinner
	}
}

func WithShaPinning(shaPinning bool) UpdaterOpt {
	return func(u *Updater) {
		u.pinImageSha = shaPinning
	}
}
