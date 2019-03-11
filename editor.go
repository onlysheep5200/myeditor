package myeditor

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
		Get exit chan of the editor
	*/
	quitChan() chan interface{}

	/**
		Add runnint go routine count of the editor
	*/
	addGoRoutineCount(count int)
}

