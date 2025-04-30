package kite

import (
	"github.com/maxime-peim/gokite/pkg/bird"
	"github.com/maxime-peim/gokite/pkg/bird/commands"
)

type (
	BirdKiteMock     struct{}
	CommandReplyMock struct{}
)

func (c *CommandReplyMock) Parse(reply string) error {
	return nil
}

func (c *CommandReplyMock) String() string {
	return "mock"
}

func (b *BirdKiteMock) Connect() error    { return nil }
func (b *BirdKiteMock) Disconnect() error { return nil }
func (b *BirdKiteMock) SendCommand(cmd commands.Command) (commands.CommandReply, error) {
	return &CommandReplyMock{}, nil
}

func (b *BirdKiteMock) SendRawCommand(command string) error {
	return nil
}

func (b *BirdKiteMock) ReadRawReply() (*bird.RawReply, error) {
	return &bird.RawReply{}, nil
}

func NewBirdKiteMock() *BirdKite {
	return &BirdKite{
		bird: &BirdKiteMock{},
	}
}
