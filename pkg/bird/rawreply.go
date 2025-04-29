package bird

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	replycodes "github.com/maxime-peim/gokite/pkg/bird/replycodes/generated"
	"github.com/pkg/errors"
)

const (
	MinCode = 0
	MaxCode = 10000
)

type rawReplyLine struct {
	content     []byte
	code        replycodes.ReplyCode
	termination bool
}

func (r *rawReplyLine) Code() replycodes.ReplyCode {
	return r.code
}

func (r *rawReplyLine) Content() []byte {
	return r.content
}

func (r *rawReplyLine) Type() replycodes.ReplyType {
	return r.code.Type()
}

func (r *rawReplyLine) String() string {
	return fmt.Sprintf("%s [%d] (%s)", string(r.content), r.code, r.code)
}

type RawReply struct {
	content []byte
	lines   []*rawReplyLine
}

func (r *RawReply) Complete() bool {
	return len(r.lines) > 0 && r.lines[len(r.lines)-1].termination
}

func (r *RawReply) Type() replycodes.ReplyType {
	if len(r.lines) == 0 {
		return replycodes.RunTimeErrorReplyType
	}
	return r.lines[0].Type()
}

func (r *RawReply) Errored() bool {
	switch r.Type() {
	case replycodes.RunTimeErrorReplyType,
		replycodes.ParseTimeErrorReplyType:
		return true
	default:
		return false
	}
}

func parseReplyLine(replyLineBytes []byte, lastCode replycodes.ReplyCode) (*rawReplyLine, error) {
	if len(replyLineBytes) == 0 {
		return nil, errors.New("empty reply line")
	} else if replyLineBytes[0] == '+' { // async reply
		return &rawReplyLine{content: replyLineBytes[1:]}, nil
	} else if replyLineBytes[0] == ' ' { // continuation line
		return &rawReplyLine{code: lastCode, content: replyLineBytes[1:]}, nil
	} else if len(replyLineBytes) < 3 {
		return nil, errors.Errorf("invalid reply line: %s", string(replyLineBytes))
	}

	code := replycodes.ReplyCode(0)
	_, err := fmt.Sscanf(string(replyLineBytes[:4]), "%d", &code)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse reply code: %s", string(replyLineBytes[:4]))
	}

	if code < MinCode || code >= MaxCode || (replyLineBytes[4] != ' ' && replyLineBytes[4] != '-') {
		return nil, errors.Errorf("invalid reply line: %s", string(replyLineBytes))
	}

	return &rawReplyLine{
		code:        code,
		content:     replyLineBytes[5:],
		termination: replyLineBytes[4] == ' ',
	}, nil
}

func newRawReply() *RawReply {
	return &RawReply{
		content: []byte{},
		lines:   []*rawReplyLine{},
	}
}

func (r *RawReply) parse(replyBytes []byte) ([]byte, error) {
	var errs *multierror.Error

	start := 0
	lastCode := replycodes.ReplyCode(0)
	for p := range replyBytes {
		if replyBytes[p] != '\n' {
			continue
		}
		replyLine, err := parseReplyLine(replyBytes[start:p], lastCode)
		if err != nil {
			errs = multierror.Append(errs, err)
		} else {
			r.lines = append(r.lines, replyLine)
			r.content = append(r.content, replyLine.content...)
			r.content = append(r.content, '\n')
			lastCode = replyLine.code
		}
		start = p + 1
		if replyLine.termination {
			break
		}
	}
	if len(r.lines) == 0 {
		return nil, errors.New("no reply line found")
	}
	if r.Complete() && start < len(replyBytes) {
		errs = multierror.Append(errs, errors.New("too long reply line"))
	}
	return replyBytes[start:], errs.ErrorOrNil()
}

func (r *RawReply) String() string {
	return strings.TrimRight(string(r.content), "\n")
}
