package conf

type IncludeStatement struct {
	Filename string
}

func (s *IncludeStatement) Marshal(opts *MarshallingOptions) (string, error) {
	return sAppendf("", "include \"%s\"", s.Filename), nil
}
