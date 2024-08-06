package urafvcli

import (
	"github.com/urfave/cli/v2"
)

// boolFlagBuilder is a builder for creating a BoolFlag
type boolFlagBuilder struct {
	flag *cli.BoolFlag
}

// NewBoolFlagBuilder returns a new boolFlagBuilder
func NewBoolFlagBuilder(name string) *boolFlagBuilder {
	return &boolFlagBuilder{
		flag: &cli.BoolFlag{
			Name: name,
		},
	}
}

// WithAliases sets the aliases for the flag
func (b *boolFlagBuilder) WithAliases(aliases ...string) *boolFlagBuilder {
	b.flag.Aliases = aliases
	return b
}

// WithUsage sets the usage description for the flag
func (b *boolFlagBuilder) WithUsage(usage string) *boolFlagBuilder {
	b.flag.Usage = usage
	return b
}

// WithEnvVars sets the environment variables for the flag
func (b *boolFlagBuilder) WithEnvVars(envVars ...string) *boolFlagBuilder {
	b.flag.EnvVars = envVars
	return b
}

// WithFilePath sets the file path for the flag
func (b *boolFlagBuilder) WithFilePath(filePath string) *boolFlagBuilder {
	b.flag.FilePath = filePath
	return b
}

// Required sets whether the flag is required
func (b *boolFlagBuilder) Required(required bool) *boolFlagBuilder {
	b.flag.Required = required
	return b
}

// Hidden sets whether the flag is hidden from help messages
func (b *boolFlagBuilder) Hidden(hidden bool) *boolFlagBuilder {
	b.flag.Hidden = hidden
	return b
}

// WithValue sets the default value for the flag
func (b *boolFlagBuilder) WithValue(value bool) *boolFlagBuilder {
	b.flag.Value = value
	return b
}

// WithDestination sets the destination pointer for the flag
func (b *boolFlagBuilder) WithDestination(destination *bool) *boolFlagBuilder {
	b.flag.Destination = destination
	return b
}

// Build returns the constructed BoolFlag
func (b *boolFlagBuilder) Build() *cli.BoolFlag {
	return b.flag
}
