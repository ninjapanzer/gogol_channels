package renderer

import "github.com/gbin/goncurses"

type Renderer interface {
	Beep()
	Draw(string)
	DrawAt(int, int, string)
	Dimensions() (int, int)
	Start()
	End()
	Refresh()
	BufferUpdate()
	Clear()
	GetChar() goncurses.Key
}
