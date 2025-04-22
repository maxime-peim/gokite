package commands

import (
	"regexp"

	"github.com/pkg/errors"
)

var readyPattern = regexp.MustCompile(`BIRD ([\d\.]+) ready.`)

type ReadyReply struct {
	Version string
}

func (r *ReadyReply) Parse(replyBytes []byte) error {
	matches := readyPattern.FindStringSubmatch(string(replyBytes))
	if len(matches) != 2 {
		return errors.Errorf("invalid version format: %s", string(replyBytes))
	}
	r.Version = string(matches[1])
	return nil
}
