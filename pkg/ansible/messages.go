package ansible

import (
	"strings"
	"sync"
)

// MessageLog holds messages during module run
type MessageLog struct {
	messages []string
	mux      sync.Mutex
}

// Add appends a message to the log
func (l *MessageLog) Add(msg string) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.messages = append(l.messages, msg)
}

func (l *MessageLog) String() string {
	l.mux.Lock()
	defer l.mux.Unlock()
	return strings.Join(l.messages, ", ")
}
