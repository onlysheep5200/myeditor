package myeditor

import (
	"time"
)

//Event...
type Event struct {
	EventType int
	Data      map[string]interface{}
	EventLoop EventLoop
}

//event type
const (
	EventNormalKey  = iota
	EventSpecialKey
	EventResize
	EventTick
	EventUnknown    = 9999
	//TODO: mouse event capture
)

const (
	EventDataFieldRune      = "rune"
	EventDataFieldKey       = "key"
	EventDataFieldWidth     = "width"
	EventDataFieldHeight    = "height"
	EventDataFieldTimestamp = "timestamp"
)

/**
	event listener
*/
type EventListener interface {
	onEvent(eve Event) error
}

type EventListenerConfig struct {
	TargetType  int
	TargetCode  int
	Name        string
	IsScheduled bool
	Interval    time.Duration
}

/**
	event loop
**/
type EventLoop interface {
	AddEventListener(config EventListenerConfig, eventListner EventListener)

	eventChan() chan Event

	Start()

	Editor() Editor
}

type EventListenerKey struct {
	Type int
	Name string
}

type CommandContext struct {
	Key   rune
	Eve   Event
	Extra map[string]interface{}
}

type CommandResult map[string]interface{}

type Command interface {
	Handle(context CommandContext) (CommandResult, error)
	Type() int
}
