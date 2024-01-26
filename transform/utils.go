package transform

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

func wildCardToRegexp(pattern string) string {
	components := regexp.MustCompile("[*%]+").Split(pattern, -1)
	if len(components) == 1 {
		// if len is 1, there are no *'s, return exact match pattern
		return "^" + "(?i)" + pattern + "$"
	}
	var sb strings.Builder
	for i, literal := range components {

		// Replace char with .*
		if i > 0 {
			sb.WriteString(".*")
		}

		// Quote any regular expression meta characters in the literal text.
		sb.WriteString(regexp.QuoteMeta(literal))
	}
	return "^" + "(?i)" + sb.String() + "$"
}

// wildCardMatch checks if a given value matches a wildcard pattern.
//
// pattern: the wildcard pattern to match against.
// value: the value to check.
// bool: true if the value matches the pattern, false otherwise.
func wildCardMatch(pattern string, value string) bool {
	result, _ := regexp.MatchString(wildCardToRegexp(strings.TrimSpace(pattern)), value)
	return result
}

func wildCardMatchs(pattern []string, value string) bool {
	for i := 0; i < len(pattern); i++ {
		if wildCardMatch(pattern[i], value) {
			return true
		}
	}
	return false
}

// FilterData filters the data based on the given schema and table patterns.
//
// It takes a pointer to a DataDef struct, a schema pattern string, and a table pattern string as parameters.
// It returns a pointer to a modified DataDef struct.
func FilterData(data *model.DataDef, schemaPattern string, tablePattern string) *model.DataDef {
	d := &model.DataDef{
		Fixed:   data.Fixed,
		Schemas: make([]model.Schema, 0),
	}

	for i := 0; i < len(data.Schemas); i++ {
		schema := data.Schemas[i]
		isSchemaMatched := schemaPattern == "" || wildCardMatchs(strings.Split(schemaPattern, ","), schema.Name)
		if isSchemaMatched {
			tables := lo.Filter(schema.Tables, func(t model.Table, _ int) bool {
				return tablePattern == "" || wildCardMatchs(strings.Split(tablePattern, ","), t.Name)
			})
			if len(tables) > 0 {
				schema.Tables = tables
				d.Schemas = append(d.Schemas, schema)
			}
		}
	}
	return d
}

// SearchPathFiles searches for files with the specified filename in directories listed in the PATH environment variable.
//
// It takes a single parameter:
// - filename: a string representing the name of the file to search for.
//
// It returns a slice of strings and an error. The slice contains the paths of the matched files, and the error is non-nil if there was an error during the search.
func SearchPathFiles(filename string) ([]string, error) {
	// Get the PATH environment variable
	path := os.Getenv("PATH")
	// add current director to the start of path
	path = "." + string(os.PathListSeparator) + path
	// Split the PATH variable into individual directories
	dirs := strings.Split(path, string(os.PathListSeparator))

	var matches []string

	// Iterate over each directory in the PATH
	for _, dir := range dirs {
		// Get a list of files in the current directory
		files, err := filepath.Glob(filepath.Join(dir, filename))
		if err != nil {
			return nil, err
		}

		// Add the matched files to the results
		matches = append(matches, files...)
	}

	if len(matches) > 0 {
		// Files found, return the matches
		return matches, nil
	}

	// No file matches found
	return nil, tracerr.Errorf("no matching files found")
}

func Verify(data *model.DataDef) []string {
	tables := make(map[string][]model.Column)

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

	result := make([]string, 0)

	for k, columns := range tables {
		for i, column := range columns {
			if column.Name == "" {
				result = append(result, fmt.Sprintf("[TB: %s] missing column name at line %d", k, i))
			}
			if column.DataType == "" {
				result = append(result, fmt.Sprintf("[TB: %s] missing data type of the column '%s' at line %d", k, column.Name, i))
			}
			if column.ForeignKey != "" {
				// check the foreign key whether exists
				if !isFKExist(column.ForeignKey) {
					result = append(result, fmt.Sprintf("[TB: %s] [FK: %s] cannot be found at line %d", k, column.ForeignKey, i))
				}
			}
		}
	}
	return result
}
