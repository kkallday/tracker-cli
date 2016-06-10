package application

import (
	"errors"
	"flag"
	"io/ioutil"
)

type CommandLineArgs struct {
	ConfigDir string
}

type FlagParser struct{}

func NewFlagParser() FlagParser {
	return FlagParser{}
}

func (f FlagParser) Parse(args []string) (CommandLineArgs, error) {
	cmdLineFlags := CommandLineArgs{}

	flagSet := flag.NewFlagSet("command line args", flag.ContinueOnError)
	flagSet.StringVar(&cmdLineFlags.ConfigDir, "config-dir", "", "path to directory containing config.json")
	flagSet.Usage = func() {}
	flagSet.SetOutput(ioutil.Discard)

	err := flagSet.Parse(args)
	if err != nil {
		return cmdLineFlags, errors.New("missing required flag --config-dir")
	}
	return cmdLineFlags, nil
}
