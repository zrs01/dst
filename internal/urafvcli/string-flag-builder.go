package urafvcli

import (
	"github.com/urfave/cli/v2"
)

// StringFlagBuilder is a builder for creating a StringFlag
type StringFlagBuilder struct {
	flag *cli.StringFlag
}

// NewStringFlagBuilder returns a new stringFlagBuilder
func NewStringFlagBuilder(name string) *StringFlagBuilder {
	return &StringFlagBuilder{
		flag: &cli.StringFlag{
			Name: name,
		},
	}
}

// WithAliases sets the aliases for the flag
func (b *StringFlagBuilder) WithAliases(aliases ...string) *StringFlagBuilder {
	b.flag.Aliases = aliases
	return b
}

// WithUsage sets the usage description for the flag
func (b *StringFlagBuilder) WithUsage(usage string) *StringFlagBuilder {
	b.flag.Usage = usage
	return b
}

// WithEnvVars sets the environment variables for the flag
func (b *StringFlagBuilder) WithEnvVars(envVars ...string) *StringFlagBuilder {
	b.flag.EnvVars = envVars
	return b
}

// WithFilePath sets the file path for the flag
func (b *StringFlagBuilder) WithFilePath(filePath string) *StringFlagBuilder {
	b.flag.FilePath = filePath
	return b
}

// TakesFile sets whether the flag accepts a file path as input
func (b *StringFlagBuilder) TakesFile(takesFile bool) *StringFlagBuilder {
	b.flag.TakesFile = takesFile
	return b
}

// Required sets whether the flag is required
func (b *StringFlagBuilder) Required(required bool) *StringFlagBuilder {
	b.flag.Required = required
	return b
}

// Hidden sets whether the flag is hidden from help messages
func (b *StringFlagBuilder) Hidden(hidden bool) *StringFlagBuilder {
	b.flag.Hidden = hidden
	return b
}

// WithValue sets the default value for the flag
func (b *StringFlagBuilder) WithValue(value string) *StringFlagBuilder {
	b.flag.Value = value
	return b
}

// WithDestination sets the destination pointer for the flag
func (b *StringFlagBuilder) WithDestination(destination *string) *StringFlagBuilder {
	b.flag.Destination = destination
	return b
}

// Build returns the constructed StringFlag
func (b *StringFlagBuilder) Build() *cli.StringFlag {
	return b.flag
}
