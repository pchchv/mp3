package consts

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
)

type Version int

type Layer int

type SamplingFrequency int
