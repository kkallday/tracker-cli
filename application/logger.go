package application

import (
	"fmt"
	"io"
	"strings"

	"github.com/kkelani/tracker-cli/trackerapi"
)

type Logger struct {
	writer io.Writer
}

func NewLogger(writer io.Writer) Logger {
	return Logger{
		writer: writer,
	}
}

func (l Logger) LogStories(stories ...trackerapi.Story) {
	for _, story := range stories {
		if story.Story_type == "feature" {
			fmt.Fprintf(l.writer, "%-9s #%-8d Pts: %d %s\n",
				l.title(story.Story_type), story.Id, story.Estimate, story.Name)
		} else {
			fmt.Fprintf(l.writer, "%-9s #%-8d %s\n",
				l.title(story.Story_type), story.Id, story.Name)
		}
	}
}

func (l Logger) Log(message string) {
	fmt.Fprintf(l.writer, "%s\n", message)
}

func (l Logger) title(word string) string {
	firstChar := strings.ToUpper(string([]byte(word)[0]))
	theRest := string([]byte(word)[1:])
	return fmt.Sprintf("%s%s", firstChar, theRest)
}
