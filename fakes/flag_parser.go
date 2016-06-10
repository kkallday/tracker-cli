package fakes

import "github.com/kkelani/tracker-cli/application"

type FlagParser struct {
	ParseCall struct {
		Receives struct {
			Args []string
		}
		Returns struct {
			CommandLineArgs application.CommandLineArgs
		}
	}
}

func (f *FlagParser) Parse(args []string) application.CommandLineArgs {
	f.ParseCall.Receives.Args = args
	return f.ParseCall.Returns.CommandLineArgs
}
