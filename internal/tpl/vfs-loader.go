package tpl

import (
	"embed"
	"io"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
)

type VSFileSystemLoader struct {
	fs  embed.FS
	dir string
}

var _ jet.Loader = (*VSFileSystemLoader)(nil)

func NewVSFileSystemLoader(fs embed.FS, dirPath string) *VSFileSystemLoader {
	return &VSFileSystemLoader{
		fs:  fs,
		dir: filepath.FromSlash(dirPath),
	}
}

func (l *VSFileSystemLoader) Exists(templatePath string) bool {
	templatePath = filepath.Join(l.dir, filepath.FromSlash(templatePath))

	// test the file whether exist in the embed.FS or not
	if _, err := l.fs.ReadFile(templatePath); err == nil {
		return true
	}
	return false
}

func (l *VSFileSystemLoader) Open(templatePath string) (io.ReadCloser, error) {
	return l.fs.Open(filepath.Join(l.dir, filepath.FromSlash(templatePath)))
}
