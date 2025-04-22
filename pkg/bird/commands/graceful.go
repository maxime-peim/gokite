package commands

type (
	GracefulRestartCommand struct{}
	GracefulRestartReply   struct{}
)

func (c *GracefulRestartCommand) String() string {
	return "graceful restart"
}

func (c *GracefulRestartCommand) NewReply() CommandReply {
	return &GracefulRestartReply{}
}

func (c *GracefulRestartReply) Parse(reply string) error {
	return nil
}

func (c *GracefulRestartReply) String() string {
	return ""
}
