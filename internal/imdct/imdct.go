package imdct

import "math"

var imdctWinData = [4][36]float32{}

func init() {
	for i := 0; i < 36; i++ {
		imdctWinData[0][i] = float32(math.Sin(math.Pi / 36 * (float64(i) + 0.5)))
	}
	for i := 0; i < 18; i++ {
		imdctWinData[1][i] = float32(math.Sin(math.Pi / 36 * (float64(i) + 0.5)))
	}
	for i := 18; i < 24; i++ {
		imdctWinData[1][i] = 1.0
	}
	for i := 24; i < 30; i++ {
		imdctWinData[1][i] = float32(math.Sin(math.Pi / 12 * (float64(i) + 0.5 - 18.0)))
	}
	for i := 30; i < 36; i++ {
		imdctWinData[1][i] = 0.0
	}
	for i := 0; i < 12; i++ {
		imdctWinData[2][i] = float32(math.Sin(math.Pi / 12 * (float64(i) + 0.5)))
	}
	for i := 12; i < 36; i++ {
		imdctWinData[2][i] = 0.0
	}
	for i := 0; i < 6; i++ {
		imdctWinData[3][i] = 0.0
	}
	for i := 6; i < 12; i++ {
		imdctWinData[3][i] = float32(math.Sin(math.Pi / 12 * (float64(i) + 0.5 - 6.0)))
	}
	for i := 12; i < 18; i++ {
		imdctWinData[3][i] = 1.0
	}
	for i := 18; i < 36; i++ {
		imdctWinData[3][i] = float32(math.Sin(math.Pi / 36 * (float64(i) + 0.5)))
	}
}
