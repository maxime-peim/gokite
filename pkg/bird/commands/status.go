package commands

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

type (
	StatusCommand struct{}
	StatusReply   struct {
		Version           string
		RouterID          string
		Hostname          string
		CurrentServerTime time.Time
		LastReboot        time.Time
		LastReconfig      time.Time
		DaemonStatus      string
	}
)

func (c *StatusCommand) String() string {
	return "show status"
}

func (c *StatusCommand) NewReply() CommandReply {
	return &StatusReply{}
}

func (c *StatusReply) Parse(reply string) error {
	lines := strings.Split(reply, "\n")
	if len(lines) < 7 {
		return errors.Errorf("invalid status reply: %s", reply)
	}
	c.Version = strings.TrimPrefix(lines[0], "BIRD ")
	c.RouterID = strings.TrimPrefix(lines[1], "Router ID is ")
	c.Hostname = strings.TrimPrefix(lines[2], "Hostname is ")
	c.CurrentServerTime, _ = time.Parse("2006-01-02 15:04:05", strings.TrimPrefix(lines[3], "Current server time is "))
	c.LastReboot, _ = time.Parse("2006-01-02 15:04:05", strings.TrimPrefix(lines[4], "Last reboot on "))
	c.LastReconfig, _ = time.Parse("2006-01-02 15:04:05", strings.TrimPrefix(lines[5], "Last reconfiguration on "))
	c.DaemonStatus = strings.TrimPrefix(lines[6], "Daemon is ")
	return nil
}

func (c *StatusReply) String() string {
	return strings.Join([]string{
		"BIRD " + c.Version,
		"Router ID is " + c.RouterID,
		"Hostname is " + c.Hostname,
		"Current server time is " + c.CurrentServerTime.Format("2006-01-02 15:04:05"),
		"Last reboot on " + c.LastReboot.Format("2006-01-02 15:04:05"),
		"Last reconfiguration on " + c.LastReconfig.Format("2006-01-02 15:04:05"),
		"Daemon is " + c.DaemonStatus,
	}, "\n")
}
