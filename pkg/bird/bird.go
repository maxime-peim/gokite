package bird

import (
	"net"
	"strings"

	"github.com/maxime-peim/gokite/pkg/bird/commands"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ReplyBufferSize = 4096

type BirdInstance struct {
	socketPath string

	conn net.Conn
	busy bool

	remaining []byte

	log *zap.SugaredLogger
}

type birdOpts struct {
	logLevel zapcore.Level
}

func WithLogLevel(level zapcore.Level) func(*birdOpts) {
	return func(opts *birdOpts) {
		opts.logLevel = level
	}
}

type BirdOption func(*birdOpts)

func NewBirdInstance(socketPath string, optFns ...BirdOption) *BirdInstance {
	opts := &birdOpts{}
	for _, fn := range optFns {
		fn(opts)
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(opts.logLevel)
	logger, _ := logConfig.Build()
	defer logger.Sync()
	return &BirdInstance{
		socketPath: socketPath,
		busy:       true,
		log:        logger.Sugar(),
	}
}

func (b *BirdInstance) readInitConnect() error {
	reply, err := b.ReadRawReply()
	if err != nil {
		return errors.Wrap(err, "failed to read init connect reply")
	}
	b.log.Debugf("init connect reply: %s", reply.String())
	readyReply := &commands.ReadyReply{}
	if err := readyReply.Parse(reply.content); err != nil {
		return errors.Wrap(err, "failed to parse init connect reply")
	} else if reply.Errored() {
		return errors.Errorf("unexpected reply type: %s\nreply: %s", reply.Type(), reply)
	}
	b.log.Infof("Bird version: %s", readyReply.Version)
	return nil
}

func (b *BirdInstance) Connect() error {
	if b.conn != nil {
		return nil
	}
	b.log.Debug("connecting to Bird socket")
	conn, err := net.Dial("unix", b.socketPath)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Bird socket")
	}
	b.log.Debug("connected to Bird socket")
	b.conn = conn
	return b.readInitConnect()
}

func (b *BirdInstance) Disconnect() error {
	if b.conn == nil {
		return nil
	}
	b.log.Debug("closing Bird socket")
	err := b.conn.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close Bird socket")
	}
	b.log.Info("Bird socket closed")
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
	b.log.Debugf("sending command to Bird: %s", command)
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
	b.log.Debug("reading reply from Bird")
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
		b.log.Debugf("read %d bytes from Bird", n)
		remaining, err := reply.parse(replyBytes[:n])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse response from Bird")
		}
		b.busy = !reply.Complete()
		b.remaining = remaining
	}
	b.log.Debug("reply from Bird complete")
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
	b.log.Debugf("parsing reply from Bird: %T", reply)
	if err := reply.Parse(rawReply.String()); err != nil {
		return nil, errors.Wrapf(err, "failed to parse command reply: %s", rawReply.String())
	}
	return reply, nil
}
