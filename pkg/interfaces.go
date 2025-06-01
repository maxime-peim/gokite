package kite

import (
	"github.com/maxime-peim/gokite/pkg/bird/commands"
	"github.com/pkg/errors"
)

func (kite *BirdKite) Interfaces() (*commands.InterfacesReply, error) {
	reply, err := kite.bird.SendCommand(&commands.InterfacesCommand{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Interfaces command")
	}
	return reply.(*commands.InterfacesReply), nil
}

func (kite *BirdKite) InterfacesSummary() (*commands.InterfacesSummaryReply, error) {
	reply, err := kite.bird.SendCommand(&commands.InterfacesCommand{
		Summary: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Interfaces summary command")
	}
	return reply.(*commands.InterfacesSummaryReply), nil
}
