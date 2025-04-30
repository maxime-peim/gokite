package kite

import (
	"github.com/maxime-peim/gokite/pkg/bird"
	"github.com/maxime-peim/gokite/pkg/bird/commands"
)

type Bird interface {
	Connect() error
	Disconnect() error
	SendCommand(cmd commands.Command) (commands.CommandReply, error)
	SendRawCommand(command string) error
	ReadRawReply() (*bird.RawReply, error)
}

type BirdKite struct {
	bird Bird
}

func NewBirdKite(socketPath string) (*BirdKite, error) {
	bird := bird.NewBirdInstance(socketPath)
	if err := bird.Connect(); err != nil {
		return nil, err
	}
	return &BirdKite{bird: bird}, nil
}
