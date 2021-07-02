package docker

import (
	"github.com/thepwagner/action-update/actions/updateaction"
	"github.com/thepwagner/action-update/updater"
)

type Environment struct {
	updateaction.Environment
	ShaPinning bool `env:"INPUT_SHA_PINNING" envDefault:"false"`
}

func (c *Environment) NewUpdater(root string) updater.Updater {
	u := NewUpdater(root, WithShaPinning(c.ShaPinning))
	u.pathFilter = c.Ignored
	return u
}
