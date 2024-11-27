package frame

import (
	"github.com/pchchv/mp3/internal/bits"
	"github.com/pchchv/mp3/internal/frameheader"
	"github.com/pchchv/mp3/internal/maindata"
	"github.com/pchchv/mp3/internal/sideinfo"
)

type Frame struct {
	header       frameheader.FrameHeader
	sideInfo     *sideinfo.SideInfo
	mainData     *maindata.MainData
	mainDataBits *bits.Bits
	store        [2][32][18]float32
	v_vec        [2][1024]float32
}
