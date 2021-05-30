package terminal

import (
	"fmt"
	"os"
)

// UpdatingStatusWriter overwrites messages on successive writes
type UpdatingStatusWriter struct {
	lastMessageLength int
}

func (w *UpdatingStatusWriter) Printf(format string, a ...interface{}) {
	// format current message
	currentMessage := fmt.Sprintf(format, a...)
	// right-pad current message and prefix with carriage return to put at start of line
	fmt.Printf("\r%-*s", w.lastMessageLength, currentMessage)

	w.lastMessageLength = len(currentMessage)
}

func IsTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
