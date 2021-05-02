package myerrors

import (
	"fmt"
	"time"
)

func NewErrorWithTime(text string) error {
	return &errorWithTime{
		text: text,
		time: time.Now(),
	}
}

type errorWithTime struct {
	text string
	time time.Time
}

func (err *errorWithTime) Error() string {
	return fmt.Sprintf("'%s' occurred at %s", err.text, err.time.Format("02-01-2006 15:04:05"))
}
