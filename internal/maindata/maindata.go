package maindata

import (
	"fmt"
	"io"

	"github.com/pchchv/mp3/internal/bits"
	"github.com/pchchv/mp3/internal/consts"
	"github.com/pchchv/mp3/internal/frameheader"
	"github.com/pchchv/mp3/internal/sideinfo"
)

var (
	nSlen2             = initSlen() /* MPEG 2.0 slen for 'normal' mode */
	scalefacSizesMpeg2 = [3][6][4]int{
		{{6, 5, 5, 5}, {6, 5, 7, 3}, {11, 10, 0, 0},
			{7, 7, 7, 0}, {6, 6, 6, 3}, {8, 8, 5, 0}},
		{{9, 9, 9, 9}, {9, 9, 12, 6}, {18, 18, 0, 0},
			{12, 12, 12, 0}, {12, 9, 9, 6}, {15, 12, 9, 0}},
		{{6, 9, 9, 9}, {6, 9, 12, 6}, {15, 18, 0, 0},
			{6, 15, 12, 0}, {6, 12, 9, 6}, {6, 18, 9, 0}}}
)

type FullReader interface {
	ReadFull([]byte) (int, error)
}

// MainData is MPEG1 Layer 3 Main Data.
type MainData struct {
	ScalefacL [2][2][22]int      // 0-4 bits
	ScalefacS [2][2][13][3]int   // 0-4 bits
	Is        [2][2][576]float32 // Huffman coded freq. lines
}

func initSlen() (nSlen2 [512]int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			n := j + i*3
			nSlen2[n+500] = i | (j << 3) | (2 << 12) | (1 << 15)
		}
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 4; k++ {
				for l := 0; l < 4; l++ {
					n := l + k*4 + j*16 + i*80
					nSlen2[n] = i | (j << 3) | (k << 6) | (l << 9) | (0 << 12)
				}
			}
		}
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 4; k++ {
				n := k + j*4 + i*20
				nSlen2[n+400] = i | (j << 3) | (k << 6) | (1 << 12)
			}
		}
	}

	return
}

func getScaleFactorsMpeg2(m *bits.Bits, header frameheader.FrameHeader, sideInfo *sideinfo.SideInfo) (*MainData, *bits.Bits, error) {
	nch := header.NumberOfChannels()
	md := &MainData{}
	for ch := 0; ch < nch; ch++ {
		part_2_start := m.BitPos()
		numbits := 0
		slen := nSlen2[sideInfo.ScalefacCompress[0][ch]]
		sideInfo.Preflag[0][ch] = (slen >> 15) & 0x1

		n := 0
		if sideInfo.BlockType[0][ch] == 2 {
			n++
			if sideInfo.MixedBlockFlag[0][ch] != 0 {
				n++
			}
		}

		var scaleFactors []int
		d := (slen >> 12) & 0x7
		for i := 0; i < 4; i++ {
			num := slen & 0x7
			slen >>= 3
			if num > 0 {
				for j := 0; j < scalefacSizesMpeg2[n][d][i]; j++ {
					scaleFactors = append(scaleFactors, m.Bits(num))
				}
				numbits += scalefacSizesMpeg2[n][d][i] * num
			} else {
				for j := 0; j < scalefacSizesMpeg2[n][d][i]; j++ {
					scaleFactors = append(scaleFactors, 0)
				}
			}
		}

		n = (n << 1) + 1
		for i := 0; i < n; i++ {
			scaleFactors = append(scaleFactors, 0)
		}

		if len(scaleFactors) == 22 {
			for i := 0; i < 22; i++ {
				md.ScalefacL[0][ch][i] = scaleFactors[i]
			}
		} else {
			for x := 0; x < 13; x++ {
				for i := 0; i < 3; i++ {
					md.ScalefacS[0][ch][x][i] = scaleFactors[(x*3)+i]
				}
			}
		}

		// read Huffman coded data. Skip stuffing bits
		if err := readHuffman(m, header, sideInfo, md, part_2_start, 0, ch); err != nil {
			return nil, nil, err
		}
	}
	// ancillary data is stored here,but we ignore it
	return md, m, nil
}

func read(source FullReader, prev *bits.Bits, size int, offset int) (*bits.Bits, error) {
	if size > 1500 {
		return nil, fmt.Errorf("mp3: size = %d", size)
	}
	// check that there's data available from previous frames if needed
	if prev != nil && offset > prev.LenInBytes() {
		// does not exist, so decoding of this frame is skipped,
		// but it is necessary to read main_data bits from the
		// bitstream in case they are needed for decoding the next frame
		buf := make([]byte, size)
		if n, err := source.ReadFull(buf); n < size {
			if err == io.EOF {
				return nil, &consts.UnexpectedEOF{At: "maindata.Read (1)"}
			}
			return nil, err
		}
		return bits.Append(prev, buf), nil
	}

	// copy data from previous frames
	vec := []byte{}
	if prev != nil {
		vec = prev.Tail(offset)
	}

	// read the main_data from file
	buf := make([]byte, size)
	if n, err := source.ReadFull(buf); n < size {
		if err == io.EOF {
			return nil, &consts.UnexpectedEOF{At: "maindata.Read (2)"}
		}
		return nil, err
	}

	return bits.New(append(vec, buf...)), nil
}
