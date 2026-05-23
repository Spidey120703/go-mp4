package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesQuickTime(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "alac",
			src: &Alac{
				FrameLength:       4096,
				CompatibleVersion: 0,
				BitDepth:          16,
				Pb:                40,
				Mb:                10,
				Kb:                14,
				NumChannels:       2,
				MaxRun:            255,
				MaxFrameBytes:     16388,
				AvgBitRate:        1411200,
				SampleRate:        44100,
			},
			dst: &Alac{
				FrameLength:       4096,
				CompatibleVersion: 0,
				BitDepth:          16,
				Pb:                40,
				Mb:                10,
				Kb:                14,
				NumChannels:       2,
				MaxRun:            255,
				MaxFrameBytes:     16388,
				AvgBitRate:        1411200,
				SampleRate:        44100,
			},
			bin: []byte{
				0x00,             // version
				0x00, 0x00, 0x00, // flags
				0x00, 0x00, 0x10, 0x00, // frame length
				0x00,       // compatible version
				0x10,       // bit depth
				0x28,       // pb
				0x0A,       // mb
				0x0E,       // kb
				0x02,       // num channels
				0x00, 0xFF, // max run
				0x00, 0x00, 0x40, 0x04, // max frame bytes
				0x00, 0x15, 0x88, 0x80, // average bit rate
				0x00, 0x00, 0xAC, 0x44, // sample rate
			},
			str: "Version=0 Flags=0x000000 FrameLength=4096 CompatibleVersion=0x0 BitDepth=0x10 Pb=0x28 Mb=0xa Kb=0xe NumChannels=0x2 MaxRun=255 MaxFrameBytes=16388 AvgBitRate=1411200 SampleRate=44100",
			ctx: Context{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal
			buf := bytes.NewBuffer(nil)
			n, err := Marshal(buf, tc.src, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(len(tc.bin)), n)
			assert.Equal(t, tc.bin, buf.Bytes())

			// Unmarshal
			r := bytes.NewReader(tc.bin)
			n, err = Unmarshal(r, uint64(len(tc.bin)), tc.dst, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, tc.dst)
			s, err := r.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), s)

			// UnmarshalAny
			dst, n, err := UnmarshalAny(bytes.NewReader(tc.bin), tc.src.GetType(), uint64(len(tc.bin)), tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, dst)
			s, err = r.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), s)

			// Stringify
			str, err := Stringify(tc.src, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, tc.str, str)
		})
	}
}
