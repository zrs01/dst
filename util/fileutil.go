package util

import (
	"os"

	"github.com/ztrue/tracerr"
	"gopkg.in/yaml.v3"
)

// SearchPathFiles searches for files with the specified filename in directories listed in the PATH environment variable.
// func SearchPathFiles(filename string) ([]string, error) {
// 	// Get the PATH environment variable
// 	path := os.Getenv("PATH")
// 	// add current director to the start of path
// 	path = "." + string(os.PathListSeparator) + path
// 	// Split the PATH variable into individual directories
// 	dirs := strings.Split(path, string(os.PathListSeparator))

// 	var matches []string

// 	// Iterate over each directory in the PATH
// 	for _, dir := range dirs {
// 		// Get a list of files in the current directory
// 		files, err := filepath.Glob(filepath.Join(dir, filename))
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Add the matched files to the results
// 		matches = append(matches, files...)
// 	}

// 	if len(matches) > 0 {
// 		// Files found, return the matches
// 		return matches, nil
// 	}

// 	// No file matches found
// 	return nil, tracerr.Errorf("no matching files found")
// }

func UnmarshalYml[T any](fileName string, v *T) (*T, error) {
	yamlFile, err := os.ReadFile(fileName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	if err := yaml.Unmarshal(yamlFile, v); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return v, err
}
