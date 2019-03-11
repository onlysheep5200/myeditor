package myeditor

import (
	"time"
	"log"
	"github.com/nsf/termbox-go"
)

func main(){
	err := termbox.Init()
	if err != nil{
		log.Fatal(err)
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(0, 0)
	time.Sleep(60 * time.Second)
	termbox.Close()
}