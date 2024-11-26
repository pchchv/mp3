
type FullReader interface {
	ReadFull([]byte) (int, error)
}

// MainData is MPEG1 Layer 3 Main Data.
type MainData struct {
	ScalefacL [2][2][22]int      // 0-4 bits
	ScalefacS [2][2][13][3]int   // 0-4 bits
	Is        [2][2][576]float32 // Huffman coded freq. lines
}
