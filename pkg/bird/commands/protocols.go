package commands

import (
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type (
	ProtocolsCommand struct {
		All bool
	}
	ProtocolReply struct {
		Name    string
		Proto   string
		Table   string
		StateUp bool
		Since   time.Time
		Info    string
	}
	ProtocolAllReply  struct{}
	ProtocolsReply    []ProtocolReply
	ProtocolsAllReply []ProtocolsAllReply
)

var (
	protocolTimeFormat          = "15:04:05.000"
	protocolsAllReplyLinesRegex = []*regexp.Regexp{
		regexp.MustCompile(`^([a-zA-Z0-9\.\-]+)\s+(up|down)\s+\(index=(\d+)\)$`), // Interface line
		regexp.MustCompile(`^\s+((?:\S+\s*)+)\s+MTU=(\d+)$`),                     // Attributes line
		regexp.MustCompile(`^\s+(\d+\.\d+\.\d+\.\d+/\d+)\s+(.+)$`),               // IPv4 address line
		regexp.MustCompile(`^\s+([0-9a-fA-F:]+/\d+)\s+(.+)$`),                    // IPv6 address line
	}
	protocolsReplyRegex = regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+(up|down)\s+(\S+)\s+(\S+)?$`)
)

func (c *ProtocolsCommand) String() string {
	cmdStr := CommandString{"show", "protocols"}
	cmdStr = cmdStr.AppendIf(c.All, "all")
	return cmdStr.String()
}

func (c *ProtocolsCommand) NewReply() CommandReply {
	return ternaryValue[CommandReply](c.All, &ProtocolsAllReply{}, &ProtocolsReply{})
}

func (c *ProtocolsReply) Parse(reply string) error {
	lines := strings.Split(reply, "\n")
	protocols := ProtocolsReply{}
	// skip first line as it is the header
	for _, line := range lines[1:] {
		matches := protocolsReplyRegex.FindStringSubmatch(line)
		if matches == nil {
			return errors.Errorf("invalid protocols line: %s", line)
		}
		since, err := time.Parse(protocolTimeFormat, matches[5])
		if err != nil {
			return errors.Wrapf(err, "cannot parse time: %s", matches[5])
		}
		info := ""
		if len(matches) > 6 {
			info = matches[6]
		}
		protocols = append(protocols, ProtocolReply{
			Name:    matches[1],
			Proto:   matches[2],
			Table:   matches[3],
			StateUp: matches[4] == "up",
			Since:   since,
			Info:    info,
		})
	}
	*c = protocols
	return nil
}

func (c *ProtocolsReply) String() string {
	s := strings.Builder{}
	for _, proto := range *c {
		s.WriteString(strings.Join([]string{proto.Name, proto.Proto, proto.Table, ternaryValue(proto.StateUp, "up", "down"), proto.Since.Format(protocolTimeFormat), proto.Info}, " ") + "\n")
	}
	return s.String()
}

func (c *ProtocolsAllReply) Parse(reply string) error {
	return nil
}

func (c *ProtocolsAllReply) String() string {
	return ""
}
