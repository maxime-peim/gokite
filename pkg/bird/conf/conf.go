package conf

import "strings"

type ConfStatement interface {
	Marshal(opts *MarshallingOptions) (string, error)
}

type ConfStatements []ConfStatement

func (c ConfStatements) Marshal(opts *MarshallingOptions) (string, error) {
	b := ""
	for _, statement := range c {
		m, err := genericMarshal(statement, opts)
		if err != nil {
			return "", err
		}
		b += m + ";\n"
	}
	return b, nil
}

type MarshallingOptions struct {
	Indentation string
}

var DefaultMarshallingOptions = &MarshallingOptions{
	Indentation: "  ",
}

type marshallingOptsFunc func(*MarshallingOptions)

func WithIndentation(indent int) marshallingOptsFunc {
	return func(opts *MarshallingOptions) {
		opts.Indentation = strings.Repeat(" ", indent)
	}
}

type BirdConf struct {
	Statements ConfStatements
	Filename   string
}

func (b *BirdConf) Marshal(optsFn ...marshallingOptsFunc) (conf string, err error) {
	opts := *DefaultMarshallingOptions
	for _, opt := range optsFn {
		opt(&opts)
	}
	return genericMarshal(b.Statements, &opts)
}
