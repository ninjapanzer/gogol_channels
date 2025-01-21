package renderer

import "github.com/gbin/goncurses"

type Renderer interface {
	Beep()
	Draw(string)
	DrawAt(int, int, string)
	Dimensions() (y int, x int)
	Start()
	End()
	Refresh()
	BufferUpdate()
	Clear()
	GetChar() goncurses.Key
}
