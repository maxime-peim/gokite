package kite

import (
	"github.com/maxime-peim/gokite/pkg/bird/commands"
	"github.com/pkg/errors"
)

func (kite *BirdKite) Down() error {
	_, err := kite.bird.SendCommand(&commands.DownCommand{})
	if err != nil {
		return errors.Wrap(err, "failed to send Down command")
	}
	return nil
}
