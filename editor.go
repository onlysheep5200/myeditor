package myeditor

import "sync"

const(
	ModReadWrite = iota
	ModReadOnly
)


type EditorConfig struct{
	Path string
	Mod byte
}

type Editor interface{
	
	/**
		Init editor
	*/
	Init(conf EditorConfig) error

	/**
		Start run editor

	*/
	Run() (EventLoop, error)

	/**
		Exit editor
	*/
	Exit() error

	/**
		Current Size of editor
	 */
	Size() (x,y int)

	/**
		Cursor position of this editor
	 */
	 Cursor() Cursor

	 LineOfCursor() int

	 setLineOfCursor(lineNum int)

	/** 
		Get exit chan of the editor
	*/
	quitChan() chan interface{}

	/**
		Add runnint go routine count of the editor
	*/
	addGoRoutineCount(count int)

	/**
		Get all lines of editor
	 */
	 lines() []line
}

const(
	TokenTypeSpace = iota
	TokenTypeNormal
	TokenTypeKeyword
)

type token struct {

	raw []rune //raw runes of the token

	flag tokenflag //token flag

	tokenType byte

	rendered []rune //runes after rendered

	startX int

	startY int

	endX int

	endY int
}

type tokenflag struct{
	color uint
}

type line struct{
	tokens []token //tokens of current line

	displayed bool //whether current line display in window

	startX int	//top left x

	startY int	//top left y

	endX int //bottom right x

	endY int //bottom right y

	lineLock sync.RWMutex //a read lock for current line

}

type Cursor struct{
	X int
	Y int
}

