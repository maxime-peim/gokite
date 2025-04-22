package commands

import (
	"fmt"
	"strings"
)

type CommandString []string

func (c CommandString) String() string {
	return strings.Join(c, " ")
}

func (c CommandString) Append(args ...string) CommandString {
	return append(c, args...)
}

func (c CommandString) Appendf(format string, args ...any) CommandString {
	return append(c, fmt.Sprintf(format, args...))
}

func (c CommandString) AppendIf(condition bool, args ...string) CommandString {
	if condition {
		return append(c, args...)
	}
	return c
}

func (c CommandString) AppendfIf(condition bool, format string, args ...any) CommandString {
	if condition {
		return append(c, fmt.Sprintf(format, args...))
	}
	return c
}

func (c CommandString) AppendIfNotEmpty(args ...string) CommandString {
	for _, arg := range args {
		if arg != "" {
			c = append(c, arg)
		}
	}
	return c
}

func (c CommandString) AppendValue(name string, value any) CommandString {
	c = append(c, name, fmt.Sprintf("%v", value))
	return c
}

func (c CommandString) AppendValueIf(condition bool, name string, value any) CommandString {
	if condition {
		c = append(c, name, fmt.Sprintf("%v", value))
	}
	return c
}
