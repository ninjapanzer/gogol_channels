package renderer

import (
	"github.com/gbin/goncurses"
	glog "gogol2/log"
)

type ShellRenderer struct {
	screen  *goncurses.Window
	wrapper *goncurses.Window
	Display *goncurses.Window
	Padding int
}

func NewShellRenderer(padding int) Renderer {
	screen, err := goncurses.Init()
	if err != nil {
		panic(err)
	}
	goncurses.Echo(false)
	goncurses.Cursor(0)
	//screen.Keypad(true)
	//if goncurses.MouseOk() {
	//	glog.GetLogger().Warn("Mouse support not detected.")
	//}
	//goncurses.MouseMask(goncurses.M_B1_PRESSED, nil)

	s := ShellRenderer{
		screen:  screen,
		Padding: padding,
	}
	s.screen.Timeout(0)
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
	glog.GetLogger().Info("Starting Window", "height", y, "width", x)
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

// Returns y,x
func (s *ShellRenderer) Dimensions() (y int, x int) {
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
	s.screen.Refresh()
}

func (s *ShellRenderer) BufferUpdate() {
	s.Display.NoutRefresh()
}

func (s *ShellRenderer) Clear() {
	s.Display.Clear()
}

func (s *ShellRenderer) GetChar() goncurses.Key {
	return s.screen.GetChar()
}
