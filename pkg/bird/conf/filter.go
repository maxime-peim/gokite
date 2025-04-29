package conf

type FilterStatements struct {
	ConfStatements
}

type Filter struct {
	Name       string
	Statements FilterStatements
}
