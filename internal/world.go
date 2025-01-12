package internal

import (
	"gogol2/renderer"
)

type WorldRenderer[T ChannelCell] struct {
	r     renderer.Renderer
	cells [][]T
}

func NewLife(state bool) ChannelCell {
	return ChannelCell{state: state}
}

func NewWorld(r renderer.Renderer) *WorldRenderer[ChannelCell] {
	x, y := r.Dimensions()
	cells := make([][]ChannelCell, x)
	for i := range cells {
		cells[i] = make([]ChannelCell, y)
	}
	return &WorldRenderer[ChannelCell]{
		r:     r,
		cells: cells,
	}
}

func (w *WorldRenderer[ChannelCell]) Refresh() {
	//for y, row := range w.world.Cells() {
	//	for x, cell := range row {
	//		if cell.State() {
	//			w.r.DrawAt(y, x, "0")
	//		} else {
	//			w.r.DrawAt(y, x, "-")
	//		}
	//	}
	//}
	//w.r.Refresh()
}

func (w *WorldRenderer[ChannelCell]) ComputeState() {
}

func (w *WorldRenderer[ChannelCell]) Cells() [][]ChannelCell {
	return w.cells
}

func (w *WorldRenderer[ChannelCell]) Bootstrap() {
	w.cells[0][1] = ChannelCell{state: true}
	w.cells[1][2] = ChannelCell{state: true}
	w.cells[2][0] = ChannelCell{state: true}
	w.cells[2][1] = ChannelCell{state: true}
	w.cells[2][2] = ChannelCell{state: true}
	w.cells[2][3] = ChannelCell{state: true}
	w.cells[2][4] = ChannelCell{state: true}
	w.cells[4][1] = ChannelCell{state: true}
	w.cells[4][2] = ChannelCell{state: true}
	w.cells[4][0] = ChannelCell{state: true}
	w.cells[4][1] = ChannelCell{state: true}
	w.cells[4][2] = ChannelCell{state: true}
	w.cells[4][3] = ChannelCell{state: true}
	w.cells[4][4] = ChannelCell{state: true}
	w.cells[5][1] = ChannelCell{state: true}
	w.cells[5][2] = ChannelCell{state: true}
	w.cells[5][0] = ChannelCell{state: true}
	w.cells[5][1] = ChannelCell{state: true}
	w.cells[5][2] = ChannelCell{state: true}
	w.cells[5][3] = ChannelCell{state: true}
	w.cells[5][4] = ChannelCell{state: true}
}
