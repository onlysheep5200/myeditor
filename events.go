package myeditor

import (
	"time"
)

//Event...
type Event struct{
	EventType int
	EventCode int
	Extra interface{}
}

//event type
const (
	EventNormalKey = iota
	EventSpecialKey 
)

/**
	event listener
*/
type EventListener interface{
	onEvent(eve Event) error
}

type EventListenerConfig struct{
	TargetType byte
	TargetCode int
	Name string
	IsScheduled bool
	Interval time.Duration
}

/**
	event loop
**/
type EventLoop interface{

	AddEventListener(config EventListenerConfig, eventListner EventListener)

	eventChan() chan Event
}

type EventListenerKey struct{
	Type byte
	Code int
	Name string
}


type CommondContext struct{
	Key rune
	Eve Event
	Extra map[string]interface{}
}

type CommondResult map[string]interface{}

type Commond interface{
	Handle(context CommondContext) (CommondResult, error)
	Type() int
}



