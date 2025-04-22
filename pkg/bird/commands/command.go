package commands

type Command interface {
	String() string
	NewReply() CommandReply
}

type CommandReply interface {
	String() string
	Parse(string) error
}
