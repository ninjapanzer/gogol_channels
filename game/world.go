package game

type World[T Life] interface {
	Cells() [][]T
	ComputeState()
	Bootstrap()
}
