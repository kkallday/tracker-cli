package fakes

import "github.com/kkelani/tracker-cli/trackerapi"

type Logger struct {
	LogCall struct {
		CallCount int
		Receives  struct {
			Message string
		}
	}
	LogStoriesCall struct {
		CallCount int
		Receives  struct {
			Stories []trackerapi.Story
		}
		Returns struct {
			Error error
		}
	}
}

func (l *Logger) LogStories(stories ...trackerapi.Story) error {
	l.LogStoriesCall.CallCount++
	l.LogStoriesCall.Receives.Stories = stories
	return l.LogStoriesCall.Returns.Error
}

func (l *Logger) Log(message string) {
	l.LogCall.CallCount++
	l.LogCall.Receives.Message = message
}
