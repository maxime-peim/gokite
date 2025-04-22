package kite

import (
	"github.com/maxime-peim/gokite/pkg/bird/commands"
	"github.com/pkg/errors"
)

func (kite *BirdKite) Status() (*commands.StatusReply, error) {
	reply, err := kite.bird.SendCommand(&commands.StatusCommand{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Status command")
	}
	return reply.(*commands.StatusReply), nil
}
