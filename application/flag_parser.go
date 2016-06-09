package application

import (
	"errors"
	"flag"
	"io/ioutil"
)

type CommandLineConfig struct {
	ConfigDir string
}

type FlagParser struct{}

func NewFlagParser() FlagParser {
	return FlagParser{}
}

func (f FlagParser) Parse(args []string) (CommandLineConfig, error) {
	cmdLineFlags := CommandLineConfig{}

	flagSet := flag.NewFlagSet("command line flags", flag.ContinueOnError)
	flagSet.StringVar(&cmdLineFlags.ConfigDir, "config-dir", "", "path to directory containing config.json")
	flagSet.Usage = func() {}
	flagSet.SetOutput(ioutil.Discard)

	err := flagSet.Parse(args)
	if err != nil {
		return cmdLineFlags, errors.New("missing required flag --config-dir")
	}
	return cmdLineFlags, nil
}
