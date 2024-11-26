package sideinfo

import (
	"fmt"
	"io"

	"github.com/pchchv/mp3/internal/bits"
	"github.com/pchchv/mp3/internal/consts"
	"github.com/pchchv/mp3/internal/frameheader"
)

var sideInfoBitsToRead = [2][4]int{
	{ // MPEG 1
		9, 5, 3, 4,
	},
	{ // MPEG 2
		8, 1, 2, 9,
	},
}

type FullReader interface {
	ReadFull([]byte) (int, error)
}

// SideInfo is MPEG1 Layer 3 Side Information.
// [2][2] means [gr][ch].
type SideInfo struct {
	MainDataBegin     int          // 9 bits
	PrivateBits       int          // 3 bits in mono, 5 in stereo
	Scfsi             [2][4]int    // 1 bit
	Part2_3Length     [2][2]int    // 12 bits
	BigValues         [2][2]int    // 9 bits
	GlobalGain        [2][2]int    // 8 bits
	ScalefacCompress  [2][2]int    // 4 bits
	WinSwitchFlag     [2][2]int    // 1 bit
	BlockType         [2][2]int    // 2 bits
	MixedBlockFlag    [2][2]int    // 1 bit
	TableSelect       [2][2][3]int // 5 bits
	SubblockGain      [2][2][3]int // 3 bits
	Region0Count      [2][2]int    // 4 bits
	Region1Count      [2][2]int    // 3 bits
	Preflag           [2][2]int    // 1 bit
	ScalefacScale     [2][2]int    // 1 bit
	Count1TableSelect [2][2]int    // 1 bit
	Count1            [2][2]int    // Not in file, calc by huffman decoder
}

func Read(source FullReader, header frameheader.FrameHeader) (*SideInfo, error) {
	nch := header.NumberOfChannels()
	framesize, err := header.FrameSize()
	if err != nil {
		return nil, err
	} else if framesize > 2000 {
		return nil, fmt.Errorf("mp3: framesize = %d\n", framesize)
	}

	sideinfo_size := header.SideInfoSize()

	// main data size is the rest of the frame,including ancillary data
	main_data_size := framesize - sideinfo_size - 4 // sync+header
	// CRC is 2 bytes
	if header.ProtectionBit() == 0 {
		main_data_size -= 2
	}

	// Read sideinfo from bitstream into buffer used by Bits()
	buf := make([]byte, sideinfo_size)
	n, err := source.ReadFull(buf)
	if n < sideinfo_size {
		if err == io.EOF {
			return nil, &consts.UnexpectedEOF{At: "sideinfo.Read"}
		}
		return nil, fmt.Errorf("mp3: couldn't read sideinfo %d bytes: %v", sideinfo_size, err)
	}
	s := bits.New(buf)

	mpeg1Frame := header.LowSamplingFrequency() == 0
	bitsToRead := sideInfoBitsToRead[header.LowSamplingFrequency()]

	// parse audio data
	// pointer to where we should start reading main data
	si := &SideInfo{}
	si.MainDataBegin = s.Bits(bitsToRead[0])
	// get private bits. Not used for anything
	if header.Mode() == consts.ModeSingleChannel {
		si.PrivateBits = s.Bits(bitsToRead[1])
	} else {
		si.PrivateBits = s.Bits(bitsToRead[2])
	}

	if mpeg1Frame {
		// get scale factor selection information
		for ch := 0; ch < nch; ch++ {
			for scfsi_band := 0; scfsi_band < 4; scfsi_band++ {
				si.Scfsi[ch][scfsi_band] = s.Bits(1)
			}
		}
	}
	// get the rest of the side information
	for gr := 0; gr < header.Granules(); gr++ {
		for ch := 0; ch < nch; ch++ {
			si.Part2_3Length[gr][ch] = s.Bits(12)
			si.BigValues[gr][ch] = s.Bits(9)
			si.GlobalGain[gr][ch] = s.Bits(8)
			si.ScalefacCompress[gr][ch] = s.Bits(bitsToRead[3])
			si.WinSwitchFlag[gr][ch] = s.Bits(1)
			if si.WinSwitchFlag[gr][ch] == 1 {
				si.BlockType[gr][ch] = s.Bits(2)
				si.MixedBlockFlag[gr][ch] = s.Bits(1)
				for region := 0; region < 2; region++ {
					si.TableSelect[gr][ch][region] = s.Bits(5)
				}

				for window := 0; window < 3; window++ {
					si.SubblockGain[gr][ch][window] = s.Bits(3)
				}

				if si.BlockType[gr][ch] == 2 && si.MixedBlockFlag[gr][ch] == 0 {
					si.Region0Count[gr][ch] = 8 // Implicit
				} else {
					si.Region0Count[gr][ch] = 7 // Implicit
				}
				// implicit
				si.Region1Count[gr][ch] = 20 - si.Region0Count[gr][ch]
			} else {
				for region := 0; region < 3; region++ {
					si.TableSelect[gr][ch][region] = s.Bits(5)
				}
				si.Region0Count[gr][ch] = s.Bits(4)
				si.Region1Count[gr][ch] = s.Bits(3)
				si.BlockType[gr][ch] = 0 // Implicit
				if !mpeg1Frame {
					si.MixedBlockFlag[0][ch] = 0
				}
			}

			if mpeg1Frame {
				si.Preflag[gr][ch] = s.Bits(1)
			}

			si.ScalefacScale[gr][ch] = s.Bits(1)
			si.Count1TableSelect[gr][ch] = s.Bits(1)
		}
	}

	return si, nil
}
