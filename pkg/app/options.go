package app

import (
	cliflag "github.com/nico612/voyage/pkg/cli/flag"
)

type CliOptions interface {
	Flags() (nfs cliflag.NamedFlagSets) // (命名) 标志集
	Validate() []error                  // 验证标志的方法
}

// ConfigurableOptions abstracts configuration options for reading parameters
// from a configuration file.
type ConfigurableOptions interface {
	// ApplyFlags parsing parameters from the command line or configuration file
	// to the options instance.
	ApplyFlags() []error
}

// CompleteableOptions abstracts options which can be completed.
type CompleteableOptions interface {
	Complete() error
}

// PrintableOptions abstracts options which can be printed.
type PrintableOptions interface {
	String() string
}
