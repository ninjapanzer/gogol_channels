package renderer

import (
	"github.com/gbin/goncurses"
)

type ShellRenderer struct {
	screen  *goncurses.Window
	wrapper *goncurses.Window
	Display *goncurses.Window
	Padding int
}

func NewShellRenderer() Renderer {
	screen, err := goncurses.Init()
	if err != nil {
		panic(err)
	}

	s := ShellRenderer{
		screen:  screen,
		Padding: 2,
	}
	s.Start()
	return &s
}

func (s *ShellRenderer) Start() {
	y, x := s.Dimensions()
	if s.wrapper != nil {
		if err := s.wrapper.Delete(); err != nil {
			panic(err)
		}
	}
	w, err := goncurses.NewWindow(y, x, 0, 0)
	if err != nil {
		panic(err)
	}
	if err := w.Box('|', '-'); err != nil {
		panic(err)
	}
	s.wrapper = w

	if s.Display != nil {
		if err := s.Display.Delete(); err != nil {
			panic(err)
		}
	}

	d, err := goncurses.NewWindow(y-(s.Padding*2), x-(s.Padding*2), s.Padding, s.Padding)
	if err != nil {
		panic(err)
	}
	s.Display = d

	s.wrapper.Refresh()
	s.Display.Refresh()
}

func (s *ShellRenderer) End() {
	goncurses.End()
}

func (s *ShellRenderer) Dimensions() (int, int) {
	return s.screen.MaxYX()
}

func (s *ShellRenderer) Draw(str string) {
	s.Display.Println(str)
}

func (s *ShellRenderer) DrawAt(y, x int, ach string) {
	s.Display.MovePrint(y, x, ach)
}

func (s *ShellRenderer) Beep() {
	goncurses.Beep()
}

func (s *ShellRenderer) Refresh() {
	s.Display.Refresh()
}

func (s *ShellRenderer) BufferUpdate() {
	s.Display.NoutRefresh()
}

func (s *ShellRenderer) Clear() {
	s.Display.Clear()
}
