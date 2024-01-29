package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ztrue/tracerr"
)

func WildCardToRegexp(pattern string) string {
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

// WildCardMatch checks if a given value matches a wildcard pattern.
//
// pattern: the wildcard pattern to match against.
// value: the value to check.
// bool: true if the value matches the pattern, false otherwise.
func WildCardMatch(pattern string, value string) bool {
	result, _ := regexp.MatchString(WildCardToRegexp(strings.TrimSpace(pattern)), value)
	return result
}

func WildCardMatchs(pattern []string, value string) bool {
	for i := 0; i < len(pattern); i++ {
		if WildCardMatch(pattern[i], value) {
			return true
		}
	}
	return false
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
