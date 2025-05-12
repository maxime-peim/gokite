package commands

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type (
	ConfigureCommand struct {
		Soft    bool
		File    string
		Timeout int
	}
	ConfigureReply struct {
		File string
	}
)

func (c *ConfigureCommand) String() string {
	cmdStr := CommandString{"configure"}
	cmdStr = cmdStr.AppendIf(c.Soft, "soft")
	cmdStr = cmdStr.AppendIf(c.File != "", quote(c.File))
	cmdStr = cmdStr.AppendValueIf(c.Timeout != 0, "timeout", c.Timeout)
	return cmdStr.String()
}

func (c *ConfigureCommand) NewReply() CommandReply {
	return &ConfigureReply{}
}

func (c *ConfigureReply) Parse(reply string) error {
	lines := strings.Split(reply, "\n")
	if len(lines) == 0 {
		return errors.Errorf("empty reply")
	}

	fileLine := strings.TrimSpace(lines[0])
	if strings.HasPrefix(fileLine, "Reading configuration from ") {
		c.File = strings.TrimPrefix(fileLine, "Reading configuration from ")
	} else {
		return errors.Errorf("unexpected reply format: %s", fileLine)
	}

	if lines[len(lines)-1] != "Reconfigured" {
		return errors.Errorf("unexpected last line format: %s", lines[len(lines)-1])
	}
	return nil
}

func (c *ConfigureReply) String() string {
	if c.File != "" {
		return fmt.Sprintf("Reconfigured from %s", c.File)
	}
	return "Reconfigured"
}

type (
	ConfigureConfirmCommand struct{}
	ConfigureConfirmReply   struct{}
)

func (c *ConfigureConfirmCommand) String() string {
	return "configure confirm"
}

func (c *ConfigureConfirmCommand) NewReply() CommandReply {
	return &ConfigureConfirmReply{}
}

func (c *ConfigureConfirmReply) Parse(reply string) error {
	return nil
}

func (c *ConfigureConfirmReply) String() string {
	return ""
}

type (
	ConfigureUndoCommand struct{}
	ConfigureUndoReply   struct{}
)

func (c *ConfigureUndoCommand) String() string {
	return "configure undo"
}

func (c *ConfigureUndoCommand) NewReply() CommandReply {
	return &ConfigureUndoReply{}
}

func (c *ConfigureUndoReply) Parse(reply string) error {
	return nil
}

func (c *ConfigureUndoReply) String() string {
	return ""
}

type (
	ConfigureCheckCommand struct {
		File string
	}
	ConfigureCheckReply struct {
		Status string
	}
)

func (c *ConfigureCheckCommand) String() string {
	cmdStr := CommandString{"configure", "check"}
	cmdStr = cmdStr.AppendIf(c.File != "", quote(c.File))
	return cmdStr.String()
}

func (c *ConfigureCheckCommand) NewReply() CommandReply {
	return &ConfigureCheckReply{}
}

func (c *ConfigureCheckReply) Parse(reply string) error {
	lines := strings.Split(reply, "\n")
	if len(lines) == 0 {
		return errors.Errorf("empty reply")
	}

	c.Status = strings.TrimSpace(lines[len(lines)-1])
	return nil
}

func (c *ConfigureCheckReply) String() string {
	return c.Status
}
