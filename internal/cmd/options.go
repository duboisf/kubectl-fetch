package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

// Options contains the result of parsing
// the command line options
type Options struct {
	AllNamespaces        bool
	IncludeNonNamespaced bool
	MaxInFlight          int
	Pattern              *regexp.Regexp
}

// GetOptions returns a new Options populated with the parsed
// command line arguments provided by `commandLineArgs`.
// If `commandLineArgs` is nil, os.Args[1:] is used.
func GetOptions(commandLineArgs []string) (*Options, error) {
	options := new(Options)
	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// commandLine.BoolVar(&options.AllNamespaces, "all-namespaces", false, "Get resources accross all namespaces")
	// commandLine.BoolVar(&options.AllNamespaces, "A", false, "Alias for --all-namespaces")
	commandLine.IntVar(&options.MaxInFlight, "parallel", 10, "Parallel calls to kubectl")
	commandLine.IntVar(&options.MaxInFlight, "p", 10, "Alias for --parallel")

	commandLine.Usage = func() {
		fmt.Fprintln(os.Stderr, "USAGE: kubectl fetch [OPTIONS]... [PATTERN]")
		fmt.Fprintln(os.Stderr, "\nwhere PATTERN is an optionnal regex used to limit the kubernetes kinds that are searched. For example, specifying the pattern 'istio' will limit the results to only the resource kinds that contains 'istio' e.g. gateways.networking.istio.io\n\nOptions:")
		commandLine.PrintDefaults()
	}

	commandLine.Parse(commandLineArgs)

	if commandLine.NArg() > 1 {
		commandLine.Usage()
		return nil, errors.New("too many args supplied")
	}
	if commandLine.NArg() == 1 {
		pattern := commandLine.Args()[0]
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("could not compile regex from pattern %q: %w", pattern, err)
		}
		options.Pattern = re
	}
	return options, nil
}
