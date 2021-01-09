package docker

import (
	"github.com/thepwagner/action-update/updater"
)

type Updater struct {
	root       string
	pathFilter func(string) bool

	tags TagLister
}

var _ updater.Updater = (*Updater)(nil)

func NewUpdater(root string, opts ...UpdaterOpt) *Updater {
	u := &Updater{
		root: root,
		tags: NewRemoteTagLister(),
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
