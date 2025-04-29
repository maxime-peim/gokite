package conf

import "net/netip"

type ProtoOptions struct {
	ConfStatements
}

type DisableOption struct {
	Disabled bool
}

func (d *DisableOption) Marshal(opts *MarshallingOptions) (string, error) {
	conf := sAppendIf("", d.Disabled, "disable", "yes")
	return conf, nil
}

type RouterIDOption struct {
	RouterID netip.Addr
}

func (r *RouterIDOption) Marshal(opts *MarshallingOptions) (string, error) {
	conf := sAppendIf("", r.RouterID.IsValid(), "router-id", r.RouterID.String())
	return conf, nil
}

type InterfacePrefix struct {
	InterfaceMask string
	Prefix        netip.Prefix
	Negate        bool
}

func (i *InterfacePrefix) Marshal(opts *MarshallingOptions) (string, error) {
	conf := sAppendIf("", i.Negate, "-")
	conf = sAppend(conf, "\""+i.InterfaceMask+"\"", i.Prefix.String())
	return conf, nil
}

type InterfaceOption struct {
	InterfacePrefix []InterfacePrefix
}

func (i *InterfaceOption) Marshal(opts *MarshallingOptions) (string, error) {
	conf := "interface"
	for _, p := range i.InterfacePrefix {
		m, err := p.Marshal(opts)
		if err != nil {
			return "", err
		}
		conf += " " + m
	}
	return conf, nil
}
