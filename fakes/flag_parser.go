package fakes

import "github.com/kkelani/tracker-cli/application"

type FlagParser struct {
	ParseCall struct {
		Receives struct {
			Args []string
		}
		Returns struct {
			CommandLineConfig application.CommandLineConfig
		}
	}
}

func (f *FlagParser) Parse(args []string) application.CommandLineConfig {
	f.ParseCall.Receives.Args = args
	return f.ParseCall.Returns.CommandLineConfig
}
