package internal

import (
	glog "gogol2/log"
	"gogol2/renderer"
)

type ChannelWorld[T ChannelCell] struct {
	r        renderer.Renderer
	cells    [][]*ChannelCell
	initProb float64
}

func NewLife(state bool) *ChannelCell {
	return NewChannelCell(state)
}

func NewChannelWorld[T ChannelCell](r renderer.Renderer, prob float64) *ChannelWorld[T] {
	x, y := r.Dimensions()
	cells := make([][]*ChannelCell, x)
	for i := range cells {
		cells[i] = make([]*ChannelCell, y)
		for j := range cells[i] {
			cells[i][j] = NewChannelCell(false)
		}
	}
	return &ChannelWorld[T]{
		r:        r,
		cells:    cells,
		initProb: prob,
	}
}

func (w *ChannelWorld[T]) Refresh() {
	w.r.Clear()
	w.DrawWorld()
	w.r.Refresh()
}

func (w *ChannelWorld[T]) ComputeState() {
}

func (w *ChannelWorld[T]) Cells() [][]*ChannelCell {
	return w.cells
}

func (w *ChannelWorld[T]) Bootstrap() {
	w.initializeProbabilisticDistributionOfLife(w.initProb)
	w.setupNeighborhood()
}

func (w *ChannelWorld[T]) initializeProbabilisticDistributionOfLife(prob float64) {
	//rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, _ := range w.cells {
		for j, _ := range w.cells[i] {
			w.cells[i][j].SetRenderer(w.DrawCell(i, j))
			w.DrawCell(i, j)(false)
			//if rng.Float64() < prob {
			//	w.cells[i][j].SilentSetState(true)
			//}
		}
	}

	w.cells[0][0].SetState(true)
	w.cells[0][1].SetState(true)
	w.cells[1][1].SetState(true)
	w.cells[1][2].SetState(true)
	w.cells[2][2].SetState(true)
	w.cells[5][2].SetState(true)

	w.r.Refresh()
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
			// Skip Self
			if i == 0 && j == 0 {
				continue
			}
			// Skip Out of bounds
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

			glog.GetLogger().Debug("Adding Neighbor", "CX", x, "CY", y, "TX", x+i, "TH", y+j)
			cell.AddChannel(cells[x+i][y+j].BroadcastChan())
			if cells[x+i][y+j].State() {
				cell.AddNeighborState(1)
			} else {
				cell.AddNeighborState(0)
			}
		}
	}
	cell.Live()
}

func (w *ChannelWorld[T]) DrawCell(y, x int) func(bool) {
	return func(state bool) {
		if state {
			w.r.DrawAt(y, x, "0")
		} else {
			w.r.DrawAt(y, x, "-")
		}
		w.r.BufferUpdate()
	}
}

func (w *ChannelWorld[T]) DrawWorld() {
	for y, row := range w.Cells() {
		for x, cell := range row {
			if cell.State() {
				w.r.DrawAt(y, x, "0")
			} else {
				w.r.DrawAt(y, x, "-")
			}
		}
	}
}
