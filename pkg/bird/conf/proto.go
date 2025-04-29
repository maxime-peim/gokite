package conf

type ProtoType int

const (
	Rip ProtoType = iota
	Ospf
	Bgp
	Static
	Direct
	Kernel
)

func (p ProtoType) String() string {
	switch p {
	case Rip:
		return "rip"
	case Ospf:
		return "ospf"
	case Bgp:
		return "bgp"
	case Static:
		return "static"
	case Direct:
		return "direct"
	case Kernel:
		return "kernel"
	default:
		return "unknown"
	}
}

func (p ProtoType) Marshal(opts *MarshallingOptions) (string, error) {
	return string(p.String()), nil
}

type Proto struct {
	Type    ProtoType
	Name    string
	From    string
	Options ProtoOptions
}

func (p *Proto) Marshal(opts *MarshallingOptions) (string, error) {
	conf := sAppend("protocol", p.Type.String(), p.Name)
	conf = sAppendIf(conf, p.From != "", "from", p.From)
	options, err := genericMarshalBrackets(p.Options, opts)
	if err != nil {
		return "", err
	}
	conf = sAppend(conf, options)
	return conf, nil
}
