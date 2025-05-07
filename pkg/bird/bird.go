package bird

import (
	"fmt"
	"net"
	"strings"

	"github.com/maxime-peim/gokite/pkg/bird/commands"
	"github.com/pkg/errors"
)

const ReplyBufferSize = 4096

type BirdInstance struct {
	socketPath string

	conn net.Conn
	busy bool

	remaining []byte
}

func NewBirdInstance(socketPath string) *BirdInstance {
	return &BirdInstance{
		socketPath: socketPath,
		busy:       true,
	}
}

func (b *BirdInstance) readInitConnect() error {
	reply, err := b.ReadRawReply()
	if err != nil {
		return errors.Wrap(err, "failed to read init connect reply")
	}
	readyReply := &commands.ReadyReply{}
	if err := readyReply.Parse(reply.content); err != nil {
		return errors.Wrap(err, "failed to parse init connect reply")
	} else if reply.Errored() {
		return errors.Errorf("unexpected reply type: %s\nreply: %s", reply.Type(), reply)
	}
	fmt.Printf("Connected to bird (%s)\n", readyReply.Version)
	return nil
}

func (b *BirdInstance) Connect() error {
	conn, err := net.Dial("unix", b.socketPath)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Bird socket")
	}
	b.conn = conn
	return b.readInitConnect()
}

func (b *BirdInstance) Disconnect() error {
	if b.conn == nil {
		return nil
	}
	err := b.conn.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close Bird socket")
	}
	b.conn = nil
	return nil
}

func (b *BirdInstance) SendRawCommand(command string) error {
	if b.conn == nil {
		return errors.New("not connected to Bird socket")
	} else if b.busy {
		return errors.New("Bird instance is busy")
	}
	b.busy = true
	command = strings.TrimRight(command, "\n") + "\n"
	_, err := b.conn.Write([]byte(command))
	if err != nil {
		return errors.Wrap(err, "failed to send command to Bird")
	}
	return nil
}

func (b *BirdInstance) ReadRawReply() (*RawReply, error) {
	if b.conn == nil {
		return nil, errors.New("not connected to Bird socket")
	} else if !b.busy {
		return nil, errors.New("no command sent to Bird")
	}
	reply := newRawReply()
	for b.busy {
		replyBytes := make([]byte, ReplyBufferSize)
		startReplyBytes := replyBytes
		if b.remaining != nil {
			copy(replyBytes, b.remaining)
			startReplyBytes = replyBytes[len(b.remaining):]
			b.remaining = nil
		}
		n, err := b.conn.Read(startReplyBytes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response from Bird")
		}
		remaining, err := reply.parse(replyBytes[:n])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse response from Bird")
		}
		b.busy = !reply.Complete()
		b.remaining = remaining
	}
	return reply, nil
}

func (b *BirdInstance) SendCommand(cmd commands.Command) (commands.CommandReply, error) {
	if err := b.SendRawCommand(cmd.String()); err != nil {
		return nil, errors.Wrap(err, "failed to send command to Bird")
	}
	rawReply, err := b.ReadRawReply()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response from Bird")
	}
	if rawReply.Errored() {
		return nil, errors.Errorf("unexpected reply type: %s\nreply: %s", rawReply.Type(), rawReply)
	}
	reply := cmd.NewReply()
	if err := reply.Parse(rawReply.String()); err != nil {
		return nil, errors.Wrapf(err, "failed to parse command reply: %s", rawReply.String())
	}
	return reply, nil
}
