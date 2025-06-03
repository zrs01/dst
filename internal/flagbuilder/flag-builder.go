package flagbuilder

func SchemaFileFlag() *StringFlag {
	return NewStringFlag("input").WithAliases("i").WithUsage(`Input file or database connection string. Supports:
	1. YAML schema file (e.g., schemal.yml)
	2. MySQL connection string in the format:
	   mysql://[user[:cred]@][protocol[(address[:port])]]/dbname[?param1=value1&...]
	3. SQL Server connection string in the format:
	   sqlserver://user:cred@host[:port][/dbname][?param1=value1&...]`)
}

func TemplateFileFlag() *StringFlag {
	return NewStringFlag("template").WithAliases("t").WithUsage("Template file to use for generating the output.")
}

func SchemaNameFlag() *StringFlag {
	return NewStringFlag("schema").WithUsage("Schema name pattern (e.g., \"my_schema*\") with wildcard support (* or %).")
}

func TableNameFlag() *StringFlag {
	return NewStringFlag("table").WithUsage("Table name pattern (e.g., \"my_table*\") with wildcard support (* or %).")
}

func OutputFileFlag() *StringFlag {
	return NewStringFlag("output").WithAliases("o").WithUsage("Output file")
}
