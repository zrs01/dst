package ddloader

type options struct {
	schemaPattern string
	tablePattern  string
	columnPattern string
}

type Option interface {
	apply(*options)
}

type schemaPatternOption string

func (s schemaPatternOption) apply(o *options) {
	o.schemaPattern = string(s)
}

func WithSchemaPattern(pattern string) Option {
	return schemaPatternOption(pattern)
}

type tablePatternOption string

func (t tablePatternOption) apply(o *options) {
	o.tablePattern = string(t)
}

func WithTablePattern(pattern string) Option {
	return tablePatternOption(pattern)
}

type columnPatternOption string

func (c columnPatternOption) apply(o *options) {
	o.columnPattern = string(c)
}

func WithColumnPattern(pattern string) Option {
	return columnPatternOption(pattern)
}
