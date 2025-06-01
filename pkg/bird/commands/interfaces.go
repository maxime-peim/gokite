package commands

import (
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type (
	InterfacesCommand struct {
		Summary bool
	}
	InterfaceReply struct {
		Name       string
		StateUp    bool
		Index      int
		Attributes []string
		MTU        int
		Addresses  map[netip.Prefix]string
	}
	InterfaceSummaryReply struct {
		Name        string
		StateUp     bool
		IPv4Address netip.Prefix
		IPv6Address netip.Prefix
	}
	InterfacesReply        []InterfaceReply
	InterfacesSummaryReply []InterfaceSummaryReply
)

var (
	interfaceReplyLinesRegex = []*regexp.Regexp{
		regexp.MustCompile(`^([a-zA-Z0-9\.\-]+)\s+(up|down)\s+\(index=(\d+)\)$`), // Interface line
		regexp.MustCompile(`^\s+((?:\S+\s*)+)\s+MTU=(\d+)$`),                     // Attributes line
		regexp.MustCompile(`^\s+(\d+\.\d+\.\d+\.\d+/\d+)\s+(.+)$`),               // IPv4 address line
		regexp.MustCompile(`^\s+([0-9a-fA-F:]+/\d+)\s+(.+)$`),                    // IPv6 address line
	}
	interfaceSummaryReplyRegex = regexp.MustCompile(`^([a-zA-Z0-9\.\-]+)\s+(up|down)\s+(\d+\.\d+\.\d+\.\d+/\d+)?\s+([0-9a-fA-F:]+/\d+)?$`)
)

func (c *InterfacesCommand) String() string {
	cmdStr := CommandString{"show", "interfaces"}
	cmdStr = cmdStr.AppendIf(c.Summary, "summary")
	return cmdStr.String()
}

func (c *InterfacesCommand) NewReply() CommandReply {
	if c.Summary {
		return &InterfacesSummaryReply{}
	}
	return &InterfacesReply{}
}

func (c *InterfacesReply) Parse(reply string) error {
	lines := strings.Split(reply, "\n")
	interfaces := InterfacesReply{}
	for i := 0; i < len(lines); {
		intf := InterfaceReply{}
		l0Matches := interfaceReplyLinesRegex[0].FindStringSubmatch(lines[i])
		l1Matches := interfaceReplyLinesRegex[1].FindStringSubmatch(lines[i+1])
		intf.Name = strings.TrimSpace(l0Matches[1])
		intf.StateUp = l0Matches[2] == "up"
		intf.Index, _ = strconv.Atoi(l0Matches[3])
		intf.Attributes = strings.Split(l1Matches[1], " ")
		intf.MTU, _ = strconv.Atoi(l1Matches[2])

		i += 2
		intf.Addresses = make(map[netip.Prefix]string)

		if i < len(lines) {
			l2Matches := interfaceReplyLinesRegex[2].FindStringSubmatch(lines[i])
			if l2Matches != nil {
				ipv4Addr := netip.MustParsePrefix(strings.TrimSpace(l2Matches[1]))
				intf.Addresses[ipv4Addr] = strings.TrimSpace(l2Matches[2])
				i++
			}
		}
		if i < len(lines) {
			l3Matches := interfaceReplyLinesRegex[3].FindStringSubmatch(lines[i])
			if l3Matches != nil {
				ipv6Addr := netip.MustParsePrefix(strings.TrimSpace(l3Matches[1]))
				intf.Addresses[ipv6Addr] = strings.TrimSpace(l3Matches[2])
				i++
			}
		}

		interfaces = append(interfaces, intf)
	}
	*c = interfaces
	return nil
}

func (c *InterfacesReply) String() string {
	s := strings.Builder{}
	for _, intf := range *c {
		s.WriteString(intf.Name + " " + ternaryValue(intf.StateUp, "up", "down") + " (index=" + strconv.Itoa(intf.Index) + ")\n")
		s.WriteString("  " + strings.Join(intf.Attributes, " ") + " MTU=" + strconv.Itoa(intf.MTU) + "\n")
		for addr, desc := range intf.Addresses {
			s.WriteString("  " + addr.String() + " " + desc + "\n")
		}
	}
	return s.String()
}

func (c *InterfacesSummaryReply) Parse(reply string) error {
	lines := strings.Split(reply, "\n")
	interfaces := InterfacesSummaryReply{}
	// skip first line which is the header
	for _, line := range lines[1:] {
		matches := interfaceSummaryReplyRegex.FindStringSubmatch(line)
		if matches == nil {
			return errors.Errorf("invalid interface summary line: %s", line)
		}
		ipv4, _ := netip.ParsePrefix(strings.TrimSpace(matches[3]))
		ipv6, _ := netip.ParsePrefix(strings.TrimSpace(matches[4]))
		interfaces = append(interfaces, InterfaceSummaryReply{
			Name:        matches[1],
			StateUp:     matches[2] == "up",
			IPv4Address: ipv4,
			IPv6Address: ipv6,
		})
	}

	*c = interfaces
	return nil
}

func (c *InterfacesSummaryReply) String() string {
	s := strings.Builder{}
	for _, intf := range *c {
		s.WriteString(intf.Name + " " + ternaryValue(intf.StateUp, "up", "down"))
		if intf.IPv4Address.IsValid() {
			s.WriteString(" " + intf.IPv4Address.String())
		}
		if intf.IPv6Address.IsValid() {
			s.WriteString(" " + intf.IPv6Address.String())
		}
		s.WriteString("\n")
	}
	return s.String()
}
