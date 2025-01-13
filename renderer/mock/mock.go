package mock

import (
	"gogol2/renderer"
	"log/slog"
)

type Renderer struct{}

func NewMockRenderer() renderer.Renderer {
	s := &Renderer{}
	return s
}

func (s *Renderer) Start() {
	slog.Debug("Start")
}

func (s *Renderer) End() {
	slog.Debug("End")
}

func (s *Renderer) Dimensions() (int, int) {
	return 10, 10
}

func (s *Renderer) Draw(str string) {
	slog.Debug("Draw")
}

func (s *Renderer) DrawAt(y, x int, ach string) {
	slog.Debug("DrawAt")
}

func (s *Renderer) Beep() {
	slog.Debug("Beep")
}

func (s *Renderer) Refresh() {
	slog.Debug("Refreshed")
}
