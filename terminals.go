package myeditor

import (
	"sync"
	"github.com/nsf/termbox-go"
	"log"
	"sync/atomic"
	"time"
)

type TerminalEditor struct {
	eventLoop      EventLoop
	exitChan       chan interface{}
	waitGroup      sync.WaitGroup
	content        []line
	displayContent []line
	width          int
	height         int
	cursor         Cursor
	cursorLineNum  int
}

func (self *TerminalEditor) LineOfCursor() int {
	return self.cursorLineNum
}

func (self *TerminalEditor) setLineOfCursor(lineNum int) {

	self.cursorLineNum = lineNum
}

func (self *TerminalEditor) Cursor() Cursor {
	return self.cursor
}

func (self *TerminalEditor) Size() (x, y int) {
	return self.width, self.height
}

func (self *TerminalEditor) Init(editorConfig EditorConfig) error {

	//init cursor
	self.cursor = Cursor{X: 0, Y: 0}

	//init termbox
	err := self.initTermbox(editorConfig)
	if err != nil {
		return err
	}

	//init event loop
	tel := &TerminalEditorEventLoop{editor: self}
	tel.init()
	self.eventLoop = tel

	//init exit chan
	self.exitChan = make(chan interface{})

	return nil
}

func (self *TerminalEditor) initTermbox(editorConfig EditorConfig) error {
	//term box
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	//init cursor
	termbox.SetCursor(self.cursor.X, self.cursor.Y)

	return nil
}

func (self *TerminalEditor) Run() (EventLoop, error) {
	self.eventLoop.Start()
	return self.eventLoop, nil
}

func (self *TerminalEditor) Exit() error {

	//exit goroutine
	self.exitChan <- 1

	//wait for waitGroup
	self.waitGroup.Wait()

	//close term box
	termbox.Close()
	return nil
}

func (self *TerminalEditor) addGoRoutineCount(count int) {
	self.waitGroup.Add(count)
}

func (self *TerminalEditor) quitChan() chan interface{} {
	return self.exitChan
}

func (self *TerminalEditor) lines() []line {
	return self.content
}

type TerminalEventListener struct {
	config   EventListenerConfig
	listener EventListener
}

func (self TerminalEventListener) onEvent(eve Event) error {
	return self.listener.onEvent(eve)
}

/**
	A Eventloop for termial 
*/
type TerminalEditorEventLoop struct {
	eveChan chan Event

	eventListeners map[int][]TerminalEventListener

	editor Editor

	//TODO: schedule event listener is a heap order by next trigger time 
}

func (self *TerminalEditorEventLoop) Start() {
	if self.editor == nil {
		return
	}

	quitChan := self.editor.quitChan()

	var closeFlag int32 = 0

	//event from termbox
	go func() {

		self.editor.addGoRoutineCount(1)
		defer self.editor.addGoRoutineCount(-1)

		if !(termbox.IsInit) {
			close(self.eventChan())
			return
		}

		for closeFlag == 0 {
			termboxEvent := termbox.PollEvent()
			parsed := Event{
				Data:      make(map[string]interface{}),
				EventType: EventUnknown,
			}

			switch termboxEvent.Type {
			case termbox.EventKey:
				parsed.Data[EventDataFieldRune] = termboxEvent.Ch
				if termboxEvent.Key != 0 {
					parsed.Data[EventDataFieldKey] = termboxEvent.Key
					parsed.EventType = EventSpecialKey
				} else {
					parsed.EventType = EventNormalKey
				}
			case termbox.EventResize:
				parsed.EventType = EventResize
				parsed.Data[EventDataFieldWidth] = termboxEvent.Width
				parsed.Data[EventDataFieldHeight] = termboxEvent.Height
			default:
				log.Println(termboxEvent)
			}

			self.eventChan() <- parsed
		}

	}()

	//event for ticks
	go func() {
		self.editor.addGoRoutineCount(1)
		ticker := time.NewTicker(1 * time.Millisecond)
		for {
			select {
			case timeEvent := <-ticker.C:
				eve := Event{EventType: EventTick, Data: make(map[string]interface{})}
				eve.Data[EventDataFieldTimestamp] = timeEvent.Unix()
				self.eventChan() <- eve // TODO: maybe block
			case <-quitChan:
				self.editor.addGoRoutineCount(-1)
				return
			}
		}
	}()

	for {
		select {
		case eve := <-self.eventChan():
			self.dispatch(eve)

		case <-quitChan:
			atomic.CompareAndSwapInt32(&closeFlag, 0, 1)
			return
		}
	}

}

func (self *TerminalEditorEventLoop) dispatch(eve Event) {
	switch eve.EventType {
	//TODO:tick
	default:
		if listeners, ok := self.eventListeners[eve.EventType]; ok {
			for _, l := range listeners {
				l.onEvent(eve)
			}
		}
	}
}

func (self *TerminalEditorEventLoop) init() {
	self.eveChan = make(chan Event, 50)
	self.eventListeners = make(map[int][]TerminalEventListener)
}

func (self *TerminalEditorEventLoop) AddEventListener(config EventListenerConfig, eventListner EventListener) {
	if (eventListner == nil) {
		return
	}
	termListener := TerminalEventListener{config: config, listener: eventListner}
	if !config.IsScheduled {
		if listeners, ok := self.eventListeners[config.TargetType]; ok {
			self.eventListeners[config.TargetType] = append(listeners, termListener)
		} else {
			self.eventListeners[config.TargetType] = []TerminalEventListener{termListener}
		}
	} else {
		//TODO: scheduled event listeners
		return
	}
}

func (self *TerminalEditorEventLoop) eventChan() chan Event {
	return self.eveChan
}
