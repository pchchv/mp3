package frame

import (
	"github.com/pchchv/mp3/internal/bits"
	"github.com/pchchv/mp3/internal/consts"
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

func (f *Frame) reorder(gr int, ch int) {
	re := make([]float32, consts.SamplesPerGr)
	_, sfBandIndicesShort := getSfBandIndicesArray(&f.header)
	// only reorder short blocks
	if (f.sideInfo.WinSwitchFlag[gr][ch] == 1) && (f.sideInfo.BlockType[gr][ch] == 2) { // Short blocks
		// check if the first two subbands
		// (=2*18 samples = 8 long or 3 short sfb's) uses long blocks
		sfb := 0
		// 2 longbl
		// sb first
		if f.sideInfo.MixedBlockFlag[gr][ch] != 0 {
			sfb = 3
		}

		next_sfb := sfBandIndicesShort[sfb+1] * 3
		win_len := sfBandIndicesShort[sfb+1] - sfBandIndicesShort[sfb]
		i := 36
		if sfb == 0 {
			i = 0
		}

		for i < consts.SamplesPerGr {
			// check if we're into the next scalefac band
			if i == next_sfb {
				// copy reordered data back to the original vector
				j := 3 * sfBandIndicesShort[sfb]
				copy(f.mainData.Is[gr][ch][j:j+3*win_len], re[0:3*win_len])
				// check if this band is above the rzero region,if so we're done
				if i >= f.sideInfo.Count1[gr][ch] {
					return
				}

				sfb++
				next_sfb = sfBandIndicesShort[sfb+1] * 3
				win_len = sfBandIndicesShort[sfb+1] - sfBandIndicesShort[sfb]
			}

			for win := 0; win < 3; win++ { // Do the actual reordering
				for j := 0; j < win_len; j++ {
					re[j*3+win] = f.mainData.Is[gr][ch][i]
					i++
				}
			}
		}

		// copy reordered data of last band back to original vector
		j := 3 * sfBandIndicesShort[12]
		copy(f.mainData.Is[gr][ch][j:j+3*win_len], re[0:3*win_len])
	}
}

func getSfBandIndicesArray(header *frameheader.FrameHeader) ([]int, []int) {
	sfreq := header.SamplingFrequency() // Setup sampling frequency index
	lsf := header.LowSamplingFrequency()
	sfBandIndicesShort := consts.SfBandIndices[lsf][sfreq][consts.SfBandIndicesShort]
	sfBandIndicesLong := consts.SfBandIndices[lsf][sfreq][consts.SfBandIndicesLong]
	return sfBandIndicesLong, sfBandIndicesShort
}
