package sideinfo

type FullReader interface {
	ReadFull([]byte) (int, error)
}
