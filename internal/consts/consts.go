package consts

import "fmt"

const (
	Version2_5      Version = 0
	VersionReserved Version = 1
	Version2        Version = 2
	Version1        Version = 3

	LayerReserved Layer = 0
	Layer3        Layer = 1
	Layer2        Layer = 2
	Layer1        Layer = 3

	SamplingFrequencyReserved SamplingFrequency = 3

	ModeStereo        Mode = 0
	ModeJointStereo   Mode = 1
	ModeDualChannel   Mode = 2
	ModeSingleChannel Mode = 3

	SamplesPerGr  = 576
	GranulesMpeg1 = 2

	SfBandIndicesLong  = 0
	SfBandIndicesShort = 1
)

var SfBandIndices = [2][3][2][]int{
	{ // MPEG 1
		{ // Layer 3
			{0, 4, 8, 12, 16, 20, 24, 30, 36, 44, 52, 62, 74, 90, 110, 134, 162, 196, 238, 288, 342, 418, 576},
			{0, 4, 8, 12, 16, 22, 30, 40, 52, 66, 84, 106, 136, 192},
		},
		{ // Layer 2
			{0, 4, 8, 12, 16, 20, 24, 30, 36, 42, 50, 60, 72, 88, 106, 128, 156, 190, 230, 276, 330, 384, 576},
			{0, 4, 8, 12, 16, 22, 28, 38, 50, 64, 80, 100, 126, 192},
		},
		{ // Layer 1
			{0, 4, 8, 12, 16, 20, 24, 30, 36, 44, 54, 66, 82, 102, 126, 156, 194, 240, 296, 364, 448, 550, 576},
			{0, 4, 8, 12, 16, 22, 30, 42, 58, 78, 104, 138, 180, 192},
		},
	},
	{ // MPEG 2
		{ // Layer 3
			{0, 6, 12, 18, 24, 30, 36, 44, 54, 66, 80, 96, 116, 140, 168, 200, 238, 284, 336, 396, 464, 522, 576},
			{0, 4, 8, 12, 18, 24, 32, 42, 56, 74, 100, 132, 174, 192},
		},
		{ // Layer 2
			{0, 6, 12, 18, 24, 30, 36, 44, 54, 66, 80, 96, 114, 136, 162, 194, 232, 278, 332, 394, 464, 540, 576},
			{0, 4, 8, 12, 18, 26, 36, 48, 62, 80, 104, 136, 180, 192},
		},
		{ // Layer 1
			{0, 6, 12, 18, 24, 30, 36, 44, 54, 66, 80, 96, 116, 140, 168, 200, 238, 284, 336, 396, 464, 522, 576},
			{0, 4, 8, 12, 18, 26, 36, 48, 62, 80, 104, 134, 174, 192},
		},
	},
}

type Version int

type Layer int

type SamplingFrequency int

type Mode int

type UnexpectedEOF struct {
	At string
}

func (u *UnexpectedEOF) Error() string {
	return fmt.Sprintf("mp3: unexpected EOF at %s", u.At)
}
