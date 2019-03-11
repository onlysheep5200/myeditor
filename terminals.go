package myeditor

import (
	"sync"
	"github.com/nsf/termbox-go"
)

type TerminalEditor struct{
	eventLoop EventLoop
	exitChan chan interface{}
	waitGroup sync.WaitGroup
}

func(self *TerminalEditor) Init(editorConfig EditorConfig) error{
	//init termbox
	err := self.initTermbox(editorConfig)
	if err != nil{
		return err
	}

	//init event loop
	tel := &TerminalEditorEventLoop{}
	tel.editor = self
	tel.init()
	self.eventLoop = tel

	//init exit chan
	self.exitChan = make(chan interface{})

	return nil
}

func(self *TerminalEditor) initTermbox(editorConfig EditorConfig) error{
	//term box
	err := termbox.Init()
	if err != nil{
		return err
	}
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	//init cursor 
	
	return nil
}

func(self *TerminalEditor) Run() (EventLoop, error){
	//TODO: running logic 
	return self.eventLoop, nil
}

func (self *TerminalEditor) Exit() error{

	//exit goroutine
	self.exitChan <- 1

	//wait for waitGroup
	self.waitGroup.Wait()

	//close term box
	termbox.Close()
	return nil
}

func (self *TerminalEditor) addGoRoutineCount(count int){
	self.waitGroup.Add(count)
}

func (self *TerminalEditor) quitChan() chan interface{} {
	return self.exitChan
}


/**
	A Eventloop for termial 
*/
type TerminalEditorEventLoop struct{
	
	eveChan chan Event

	eventListeners map[EventListenerKey][]EventListener

	editor Editor

	//TODO: schedule event listener is a heap order by next trigger time 
}

func (self *TerminalEditorEventLoop) init(){
	self.eveChan = make(chan Event, 10)
	self.eventListeners = make(map[EventListenerKey][]EventListener)
}

func (self *TerminalEditorEventLoop) AddEventListener(config EventListenerConfig ,eventListner EventListener) {
	if(eventListner == nil){
		return
	}
	if !config.IsScheduled{
		key := EventListenerKey{
			Type : config.TargetType,
			Code : config.TargetCode,
			Name : config.Name,
		}

		if listeners,ok := self.eventListeners[key]; ok{
			self.eventListeners[key] = append(listeners, eventListner)
		}else{
			self.eventListeners[key] = []EventListener{eventListner}
		}
	}else{
		//TODO: scheduled event listeners
		return
	}
}


func (self *TerminalEditorEventLoop) eventChan() chan Event{
	return self.eveChan
}

