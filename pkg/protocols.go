package kite

import (
	"github.com/maxime-peim/gokite/pkg/bird/commands"
	"github.com/pkg/errors"
)

func (kite *BirdKite) Protocols() (*commands.ProtocolsReply, error) {
	reply, err := kite.bird.SendCommand(&commands.ProtocolsCommand{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Interfaces command")
	}
	return reply.(*commands.ProtocolsReply), nil
}
