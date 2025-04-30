package ddloader

func WithSchemaPattern(pattern string) func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.schemaPattern = pattern
	}
}

func WithTablePattern(pattern string) func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.tablePattern = pattern
	}
}

func WithColumnPattern(pattern string) func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.columnPattern = pattern
	}
}

func WithUpdateReferenceTables() func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.updateReferenceTables = true
	}
}

func WithExpandFixColumns() func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.expandFixColumns = true
	}
}
