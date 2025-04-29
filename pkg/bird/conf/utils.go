package conf

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func sAppend(b string, s ...string) string {
	toJoin := []string{}
	if len(b) > 0 {
		toJoin = append(toJoin, b)
	}
	for _, c := range s {
		if len(c) > 0 {
			toJoin = append(toJoin, c)
		}
	}
	return strings.Join(toJoin, " ")
}

func sAppendf(b string, format string, args ...any) string {
	return sAppend(b, fmt.Sprintf(format, args...))
}

func sAppendIf(b string, condition bool, s ...string) string {
	if condition {
		return sAppend(b, s...)
	}
	return b
}

func sAppendfIf(b string, condition bool, format string, args ...any) string {
	if condition {
		return sAppendf(b, format, args...)
	}
	return b
}

func sPrepend(b string, s ...string) string {
	prepend := sAppend("", s...)
	if len(b) == 0 {
		return prepend
	}
	if len(prepend) > 0 {
		prepend += " "
	}
	return prepend + b
}

func sIndent(b string, indent string) string {
	if len(b) == 0 {
		return ""
	} else if len(indent) == 0 {
		return b
	}
	lines := strings.Split(b, "\n")
	result := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			line = indent + line
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func genericMarshal(stAny any, opts *MarshallingOptions) (string, error) {
	if marshaller, ok := stAny.(ConfStatement); ok {
		return marshaller.Marshal(opts)
	}
	return "", errors.Errorf("unsupported type for marshalling: %T", stAny)
}

func genericMarshalBrackets(stAny any, opts *MarshallingOptions) (string, error) {
	middle, err := genericMarshal(stAny, opts)
	if err != nil {
		return "", err
	}
	indentation := ""
	if opts != nil {
		indentation = opts.Indentation
	}
	middle = sIndent(middle, indentation)
	if len(middle) == 0 {
		return "{}", nil
	}
	return "{\n" + middle + "}", nil
}
