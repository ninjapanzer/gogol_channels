package internal

import (
	"gogol2/renderer"
	"log/slog"
)

type ChannelWorld[T ChannelCell] struct {
	r     renderer.Renderer
	cells [][]*ChannelCell
}

func NewLife(state bool) *ChannelCell {
	return NewChannelCell(state)
}

func NewChannelWorld[T ChannelCell](r renderer.Renderer) *ChannelWorld[T] {
	x, y := r.Dimensions()
	cells := make([][]*ChannelCell, x)
	for i := range cells {
		cells[i] = make([]*ChannelCell, y)
		for j := range cells[i] {
			cells[i][j] = NewChannelCell(false)
		}
	}
	return &ChannelWorld[T]{
		r:     r,
		cells: cells,
	}
}

func (w *ChannelWorld[T]) Refresh() {
	for y, row := range w.Cells() {
		for x, cell := range row {
			if cell.State() {
				w.r.DrawAt(y, x, "0")
			} else {
				w.r.DrawAt(y, x, "-")
			}
		}
	}
	w.r.Refresh()
}

func (w *ChannelWorld[T]) ComputeState() {
}

func (w *ChannelWorld[T]) Cells() [][]*ChannelCell {
	return w.cells
}

func (w *ChannelWorld[T]) Bootstrap() {
	w.cells[0][1].SilentSetState(true)
	w.cells[1][2].SilentSetState(true)
	w.cells[2][0].SilentSetState(true)
	w.cells[2][1].SilentSetState(true)
	w.cells[2][2].SilentSetState(true)
	w.cells[2][3].SilentSetState(true)
	w.cells[2][4].SilentSetState(true)
	w.cells[4][1].SilentSetState(true)
	w.cells[4][2].SilentSetState(true)
	w.cells[4][0].SilentSetState(true)
	w.cells[4][1].SilentSetState(true)
	w.cells[4][2].SilentSetState(true)
	w.cells[4][3].SilentSetState(true)
	w.cells[4][4].SilentSetState(true)
	w.cells[5][1].SilentSetState(true)
	w.cells[5][2].SilentSetState(true)
	w.cells[5][0].SilentSetState(true)
	w.cells[5][1].SilentSetState(true)
	w.cells[5][2].SilentSetState(true)
	w.cells[5][3].SilentSetState(true)
	w.cells[5][4].SilentSetState(true)
	w.setupNeighborhood()
}

func (w *ChannelWorld[T]) setupNeighborhood() {
	width := len(w.cells)
	height := len(w.cells[0])

	for i, _ := range w.cells {
		for j, _ := range w.cells[i] {
			linkNeighbors(w.cells, w.cells[i][j], i, j, width, height)
		}
	}
}

func linkNeighbors(cells [][]*ChannelCell, cell *ChannelCell, x, y, width, height int) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == x && j == y {
				continue
			}
			if x+i < 0 {
				continue
			}
			if x+i >= width {
				continue
			}
			if y+j < 0 {
				continue
			}
			if y+j >= height {
				continue
			}

			slog.Debug("Adding Neighbor", "CX", x, "CY", y, "TX", x+i, "TH", y+j)
			cell.AddChannel(cells[x+i][y+j].BroadcastChan())
		}
	}
}
