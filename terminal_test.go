package myeditor

import (
	"github.com/nsf/termbox-go"
	"testing"
)


func TestTermial(t *testing.T){
	err := termbox.Init()
	if err != nil{
		t.Error(err)
		t.Fatal()
	}
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(0, 0)
}