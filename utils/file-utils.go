package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ztrue/tracerr"
)

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
