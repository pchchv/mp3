package frame

import (
	"math"

	"github.com/pchchv/mp3/internal/bits"
	"github.com/pchchv/mp3/internal/consts"
	"github.com/pchchv/mp3/internal/frameheader"
	"github.com/pchchv/mp3/internal/maindata"
	"github.com/pchchv/mp3/internal/sideinfo"
)


var (
	powtab34 = make([]float64, 8207)
	pretab   = []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 3, 3, 3, 2, 0}
	isRatios = []float32{0.000000, 0.267949, 0.577350, 1.000000, 1.732051, 3.732051}
)

type Frame struct {
	header       frameheader.FrameHeader
	sideInfo     *sideinfo.SideInfo
	mainData     *maindata.MainData
	mainDataBits *bits.Bits
	store        [2][32][18]float32
	v_vec        [2][1024]float32
}

func (f *Frame) SamplingFrequency() (int, error) {
	return f.header.SamplingFrequencyValue()
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

func (f *Frame) stereoProcessIntensityLong(gr int, sfb int) {
	is_ratio_l := float32(0)
	is_ratio_r := float32(0)
	// check that((is_pos[sfb]=scalefac) < 7) => no intensity stereo
	if is_pos := f.mainData.ScalefacL[gr][0][sfb]; is_pos < 7 {
		sfBandIndicesLong, _ := getSfBandIndicesArray(&f.header)
		sfb_start := sfBandIndicesLong[sfb]
		sfb_stop := sfBandIndicesLong[sfb+1]
		if is_pos == 6 { // tan((6*PI)/12 = PI/2) needs special treatment!
			is_ratio_l = 1.0
			is_ratio_r = 0.0
		} else {
			is_ratio_l = isRatios[is_pos] / (1.0 + isRatios[is_pos])
			is_ratio_r = 1.0 / (1.0 + isRatios[is_pos])
		}

		// now decode all samples in this scale factor band
		for i := sfb_start; i < sfb_stop; i++ {
			f.mainData.Is[gr][0][i] *= is_ratio_l
			f.mainData.Is[gr][1][i] *= is_ratio_r
		}
	}
}

func (f *Frame) stereoProcessIntensityShort(gr int, sfb int) {
	is_ratio_l := float32(0)
	is_ratio_r := float32(0)
	_, sfBandIndicesShort := getSfBandIndicesArray(&f.header)
	// window length
	win_len := sfBandIndicesShort[sfb+1] - sfBandIndicesShort[sfb]
	// windows within the band has different scalefactors
	for win := 0; win < 3; win++ {
		// check that((is_pos[sfb]=scalefac) < 7) => no intensity stereo
		is_pos := f.mainData.ScalefacS[gr][0][sfb][win]
		if is_pos < 7 {
			sfb_start := sfBandIndicesShort[sfb]*3 + win_len*win
			sfb_stop := sfb_start + win_len
			if is_pos == 6 { // tan((6*PI)/12 = PI/2) needs special treatment!
				is_ratio_l = 1.0
				is_ratio_r = 0.0
			} else {
				is_ratio_l = isRatios[is_pos] / (1.0 + isRatios[is_pos])
				is_ratio_r = 1.0 / (1.0 + isRatios[is_pos])
			}

			// decode all samples in this scale factor band
			for i := sfb_start; i < sfb_stop; i++ {
				f.mainData.Is[gr][0][i] *= is_ratio_l
				f.mainData.Is[gr][1][i] *= is_ratio_r
			}
		}
	}
}

func (f *Frame) requantizeProcessLong(gr, ch, is_pos, sfb int) {
	sf_mult := 0.5
	if f.sideInfo.ScalefacScale[gr][ch] != 0 {
		sf_mult = 1.0
	}

	pf_x_pt := float64(f.sideInfo.Preflag[gr][ch]) * pretab[sfb]
	idx := -(sf_mult * (float64(f.mainData.ScalefacL[gr][ch][sfb]) + pf_x_pt)) + 0.25*(float64(f.sideInfo.GlobalGain[gr][ch])-210)
	tmp1 := math.Pow(2.0, idx)
	tmp2 := 0.0
	if f.mainData.Is[gr][ch][is_pos] < 0.0 {
		tmp2 = -powtab34[int(-f.mainData.Is[gr][ch][is_pos])]
	} else {
		tmp2 = powtab34[int(f.mainData.Is[gr][ch][is_pos])]
	}

	f.mainData.Is[gr][ch][is_pos] = float32(tmp1 * tmp2)
}

func (f *Frame) requantizeProcessShort(gr, ch, is_pos, sfb, win int) {
	sf_mult := 0.5
	if f.sideInfo.ScalefacScale[gr][ch] != 0 {
		sf_mult = 1.0
	}

	idx := -(sf_mult * float64(f.mainData.ScalefacS[gr][ch][sfb][win])) +
		0.25*(float64(f.sideInfo.GlobalGain[gr][ch])-210.0-
			8.0*float64(f.sideInfo.SubblockGain[gr][ch][win]))
	tmp1 := math.Pow(2.0, idx)
	tmp2 := 0.0
	if f.mainData.Is[gr][ch][is_pos] < 0 {
		tmp2 = -powtab34[int(-f.mainData.Is[gr][ch][is_pos])]
	} else {
		tmp2 = powtab34[int(f.mainData.Is[gr][ch][is_pos])]
	}

	f.mainData.Is[gr][ch][is_pos] = float32(tmp1 * tmp2)
}

func getSfBandIndicesArray(header *frameheader.FrameHeader) ([]int, []int) {
	sfreq := header.SamplingFrequency() // Setup sampling frequency index
	lsf := header.LowSamplingFrequency()
	sfBandIndicesShort := consts.SfBandIndices[lsf][sfreq][consts.SfBandIndicesShort]
	sfBandIndicesLong := consts.SfBandIndices[lsf][sfreq][consts.SfBandIndicesLong]
	return sfBandIndicesLong, sfBandIndicesShort
}
