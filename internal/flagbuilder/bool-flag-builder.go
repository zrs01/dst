package flagbuilder

import (
	"github.com/urfave/cli/v3"
)

// BoolFlag is a builder for creating a BoolFlag
type BoolFlag struct {
	flag *cli.BoolFlag
}

// NewBoolFlag returns a new boolFlagBuilder
func NewBoolFlag(name string) *BoolFlag {
	return &BoolFlag{
		flag: &cli.BoolFlag{
			Name: name,
		},
	}
}

// WithAliases sets the aliases for the flag
func (b *BoolFlag) WithAliases(aliases ...string) *BoolFlag {
	b.flag.Aliases = aliases
	return b
}

// WithUsage sets the usage description for the flag
func (b *BoolFlag) WithUsage(usage string) *BoolFlag {
	b.flag.Usage = usage
	return b
}

// WithEnvVars sets the environment variables for the flag
// func (b *BoolFlag) WithEnvVars(envVars ...string) *BoolFlag {
// 	b.flag.EnvVars = envVars
// 	return b
// }

// WithFilePath sets the file path for the flag
// func (b *BoolFlag) WithFilePath(filePath string) *BoolFlag {
// 	b.flag.FilePath = filePath
// 	return b
// }

// WithRequired sets whether the flag is required
func (b *BoolFlag) WithRequired(required bool) *BoolFlag {
	b.flag.Required = required
	return b
}

// WithHidden sets whether the flag is hidden from help messages
func (b *BoolFlag) WithHidden(hidden bool) *BoolFlag {
	b.flag.Hidden = hidden
	return b
}

// WithValue sets the default value for the flag
func (b *BoolFlag) WithValue(value bool) *BoolFlag {
	b.flag.Value = value
	return b
}

// WithDestination sets the destination pointer for the flag
func (b *BoolFlag) WithDestination(destination *bool) *BoolFlag {
	b.flag.Destination = destination
	return b
}

// Build returns the constructed BoolFlag
func (b *BoolFlag) Build() *cli.BoolFlag {
	return b.flag
}
