package flagbuilder

func InputFileFlag() *StringFlag {
	return NewStringFlag("input").WithAliases("i").WithUsage("schema definition file (yaml format)")
}

func TemplateFileFlag() *StringFlag {
	return NewStringFlag("template").WithAliases("t").WithUsage("template file to use for generating the output.")
}

func SchemaNameFlag() *StringFlag {
	return NewStringFlag("schema").WithUsage("schema name pattern (e.g., \"my_schema*\") with wildcard support (* or %).")
}

func TableNameFlag() *StringFlag {
	return NewStringFlag("table").WithUsage("table name pattern (e.g., \"my_table*\") with wildcard support (* or %).")
}

func OutputFileFlag() *StringFlag {
	return NewStringFlag("output").WithAliases("o").WithUsage("output file")
}

func DsnFlag() *StringFlag {
	return NewStringFlag("dsn").WithUsage(`database source name. Supports:
	1. MySQL connection string in the format:
	   mysql://[user[:cred]@][protocol[(address[:port])]]/dbname[?param1=value1&...]
	2. SQL Server connection string in the format:
	   sqlserver://user:cred@host[:port][/dbname][?param1=value1&...]`)
}

func CommonColumnFlag() *StringFlag {
	return NewStringFlag("ccf").WithAliases("cc").WithUsage("common column file")
}
