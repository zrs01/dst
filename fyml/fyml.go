package fyml

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/ztrue/tracerr"
	"gopkg.in/yaml.v3"
)

func ReadYaml(file string) (*DataDef, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var node yaml.Node
	if err := yaml.Unmarshal(yamlFile, &node); err != nil {
		return nil, tracerr.Wrap(err)
	}
	fmt.Printf("%+v\n", node)
	bytes, err := yaml.Marshal(&node)
	fmt.Println(string(bytes))

	var d DataDef
	if err := yaml.Unmarshal(yamlFile, &d); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &d, err
}

func WriteYaml(data *DataDef, outfile string) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return tracerr.Wrap(err)
	}
	output := string(bytes)
	// correct the names
	output = strings.ReplaceAll(output, "_column_values: ", "")
	output = strings.ReplaceAll(output, "out_columns", "columns")
	output = strings.ReplaceAll(output, "\"N\"", "N")
	output = strings.ReplaceAll(output, "\"n\"", "n")
	output = strings.ReplaceAll(output, "\"Y\"", "Y")
	output = strings.ReplaceAll(output, "\"y\"", "y")
	if outfile == "" {
		// fmt.Println(output)
	} else {
		os.WriteFile(outfile, []byte(output), fs.FileMode(0744))
	}
	return nil
}

func Verify(data *DataDef) {
	tables := make(map[string][]Column)

	// convert to map for easy searching
	for _, schema := range data.Schemas {
		for _, table := range schema.Tables {
			tables[table.Name] = table.Columns
		}
	}
	tables["fixed"] = data.Fixed

	isFKExist := func(fk string) bool {
		// fk format: table.field
		tf := strings.Split(fk, ".")
		if columns, found := tables[tf[0]]; found {
			for _, column := range columns {
				if column.Name == tf[1] {
					return true
				}
			}
		}
		return false
	}

	for k, columns := range tables {
		for i, column := range columns {
			if column.Name == "" {
				fmt.Printf("'%s', missing column name at line %d", k, i)
			}
			if column.DataType == "" {
				fmt.Printf("'%s', missing data type of the column '%s' at line %d", k, column.Name, i)
			}
			if column.ForeignKey != "" {
				// check the foreign key whether exists
				if !isFKExist(column.ForeignKey) {
					fmt.Printf("'%s', FK '%s' cannot be found at line %d\n", k, column.ForeignKey, i)
				}
			}
		}
	}
}
