package commands

type (
	DownCommand struct{}
	DownReply   struct{}
)

func (c *DownCommand) String() string {
	return "down"
}

func (c *DownCommand) NewReply() CommandReply {
	return &DownReply{}
}

func (c *DownReply) Parse(reply string) error {
	return nil
}

func (c *DownReply) String() string {
	return ""
}
