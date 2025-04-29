package conf

type TableOption struct {
	Name string
}

func (t *TableOption) Marshal(opts *MarshallingOptions) (string, error) {
	return sAppend("table", t.Name), nil
}

type FilterClause interface {
	ConfStatement
}

type ImportAll struct{}

func (i *ImportAll) Marshal(opts *MarshallingOptions) (string, error) {
	return string("all"), nil
}

type ImportNone struct{}

func (i *ImportNone) Marshal(opts *MarshallingOptions) (string, error) {
	return string("none"), nil
}

type ImportFilter struct {
	FilterName string
}

func (i *ImportFilter) Marshal(opts *MarshallingOptions) (string, error) {
	return sAppend("filter", i.FilterName), nil
}

type ImportFilterLocal struct {
	Statement FilterStatements
}

func (i *ImportFilterLocal) Marshal(opts *MarshallingOptions) (string, error) {
	conf, err := genericMarshalBrackets(i.Statement, opts)
	if err != nil {
		return "", err
	}
	return sPrepend(conf, "filter"), nil
}

type ImportOption struct {
	Filter FilterClause
}

func (i *ImportOption) Marshal(opts *MarshallingOptions) (string, error) {
	conf, err := i.Filter.Marshal(opts)
	if err != nil {
		return "", err
	}
	return sPrepend(conf, "import"), nil
}
