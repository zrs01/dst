package flagbuilder

import (
	"github.com/urfave/cli/v2"
)

// StringFlag is a builder for creating a StringFlag
type StringFlag struct {
	flag *cli.StringFlag
}

// NewStringFlag returns a new stringFlagBuilder
func NewStringFlag(name string) *StringFlag {
	return &StringFlag{
		flag: &cli.StringFlag{
			Name: name,
		},
	}
}

// WithAliases sets the aliases for the flag
func (b *StringFlag) WithAliases(aliases ...string) *StringFlag {
	b.flag.Aliases = aliases
	return b
}

// WithUsage sets the usage description for the flag
func (b *StringFlag) WithUsage(usage string) *StringFlag {
	b.flag.Usage = usage
	return b
}

// WithEnvVars sets the environment variables for the flag
func (b *StringFlag) WithEnvVars(envVars ...string) *StringFlag {
	b.flag.EnvVars = envVars
	return b
}

// WithFilePath sets the file path for the flag
func (b *StringFlag) WithFilePath(filePath string) *StringFlag {
	b.flag.FilePath = filePath
	return b
}

// TakesFile sets whether the flag accepts a file path as input
func (b *StringFlag) TakesFile(takesFile bool) *StringFlag {
	b.flag.TakesFile = takesFile
	return b
}

// Required sets whether the flag is required
func (b *StringFlag) Required(required bool) *StringFlag {
	b.flag.Required = required
	return b
}

// Hidden sets whether the flag is hidden from help messages
func (b *StringFlag) Hidden(hidden bool) *StringFlag {
	b.flag.Hidden = hidden
	return b
}

// WithValue sets the default value for the flag
func (b *StringFlag) WithValue(value string) *StringFlag {
	b.flag.Value = value
	return b
}

// WithDestination sets the destination pointer for the flag
func (b *StringFlag) WithDestination(destination *string) *StringFlag {
	b.flag.Destination = destination
	return b
}

// Build returns the constructed StringFlag
func (b *StringFlag) Build() *cli.StringFlag {
	return b.flag
}
