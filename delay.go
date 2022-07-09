package gotyphon

import (
	"time"
)

// Delay pauses the program for a given number of seconds.
func Delay(seconds time.Duration) {
	time.Sleep(seconds * time.Second)
}
