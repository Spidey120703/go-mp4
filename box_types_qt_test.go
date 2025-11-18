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
			name: "ludt",
			src:  &Ludt{},
			dst:  &Ludt{},
			bin:  nil,
			str:  ``,
		},
		{
			name: "tlou",
			src: &LoudnessEntry{
				AnyTypeBox:        AnyTypeBox{Type: BoxTypeTlou()},
				Version:           1,
				Flags:             [3]byte{},
				LoudnessBaseCount: 1,
				LoudnessBases: []LoudnessBase{
					{
						EQSetID:                0,
						DownmixID:              0,
						DRCSetID:               0,
						BsSamplePeakLevel:      647,
						BsTruePeakLevel:        644,
						MeasurementSystemForTP: 2,
						ReliabilityForTP:       3,
						MeasurementCount:       5,
						Measurements: []LoudnessMeasurement{{
							MethodDefinition:  1,
							MethodValue:       184,
							MeasurementSystem: 2,
							Reliability:       3,
						},
							{
								MethodDefinition:  3,
								MethodValue:       194,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  4,
								MethodValue:       200,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  5,
								MethodValue:       196,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  6,
								MethodValue:       54,
								MeasurementSystem: 0,
								Reliability:       0,
							},
						},
					},
				},
			},
			dst: &LoudnessEntry{
				AnyTypeBox:        AnyTypeBox{Type: BoxTypeTlou()},
				LoudnessBaseCount: 1,
				LoudnessBases: []LoudnessBase{
					{
						EQSetID:                0,
						DownmixID:              0,
						DRCSetID:               0,
						BsSamplePeakLevel:      647,
						BsTruePeakLevel:        644,
						MeasurementSystemForTP: 2,
						ReliabilityForTP:       3,
						MeasurementCount:       5,
						Measurements: []LoudnessMeasurement{{
							MethodDefinition:  1,
							MethodValue:       184,
							MeasurementSystem: 2,
							Reliability:       3,
						},
							{
								MethodDefinition:  3,
								MethodValue:       194,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  4,
								MethodValue:       200,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  5,
								MethodValue:       196,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  6,
								MethodValue:       54,
								MeasurementSystem: 0,
								Reliability:       0,
							},
						},
					},
				},
			},
			bin: []byte{
				0x01,             // version
				0x00, 0x00, 0x00, // flags
				0x01, 0x00, 0x00, 0x00, 0x28, 0x72, 0x84, 0x23, 0x05, 0x01, 0xB8, 0x23,
				0x03, 0xC2, 0x13, 0x04, 0xC8, 0x13, 0x05, 0xC4, 0x13, 0x06, 0x36, 0x00,
			},
			str: `Version=0 Flags=0x000000 LoudnessBases=[{EQSetID=0x0 DownmixID=0x0 DRCSetID=0x0 BsSamplePeakLevel=647 BsTruePeakLevel=644 MeasurementSystemForTP=0x2 ReliabilityForTP=0x3 MeasurementCount=0x5 Measurements=[{MethodDefinition=0x1 MethodValue=0xb8 MeasurementSystem=0x2 Reliability=0x3}, {MethodDefinition=0x3 MethodValue=0xc2 MeasurementSystem=0x1 Reliability=0x3}, {MethodDefinition=0x4 MethodValue=0xc8 MeasurementSystem=0x1 Reliability=0x3}, {MethodDefinition=0x5 MethodValue=0xc4 MeasurementSystem=0x1 Reliability=0x3}, {MethodDefinition=0x6 MethodValue=0x36 MeasurementSystem=0x0 Reliability=0x0}]}]`,
		},
		{
			name: "alou",
			src: &LoudnessEntry{
				AnyTypeBox:        AnyTypeBox{Type: BoxTypeAlou()},
				Version:           1,
				Flags:             [3]byte{},
				LoudnessBaseCount: 1,
				LoudnessBases: []LoudnessBase{
					{
						EQSetID:                0,
						DownmixID:              0,
						DRCSetID:               0,
						BsSamplePeakLevel:      643,
						BsTruePeakLevel:        588,
						MeasurementSystemForTP: 2,
						ReliabilityForTP:       3,
						MeasurementCount:       5,
						Measurements: []LoudnessMeasurement{
							{
								MethodDefinition:  1,
								MethodValue:       200,
								MeasurementSystem: 2,
								Reliability:       3,
							},
							{
								MethodDefinition:  3,
								MethodValue:       206,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  4,
								MethodValue:       216,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  5,
								MethodValue:       208,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  6,
								MethodValue:       54,
								MeasurementSystem: 0,
								Reliability:       0,
							},
						},
					},
				},
			},
			dst: &LoudnessEntry{
				AnyTypeBox:        AnyTypeBox{Type: BoxTypeAlou()},
				Version:           1,
				Flags:             [3]byte{},
				LoudnessBaseCount: 1,
				LoudnessBases: []LoudnessBase{
					{
						EQSetID:                0,
						DownmixID:              0,
						DRCSetID:               0,
						BsSamplePeakLevel:      643,
						BsTruePeakLevel:        588,
						MeasurementSystemForTP: 2,
						ReliabilityForTP:       3,
						MeasurementCount:       5,
						Measurements: []LoudnessMeasurement{
							{
								MethodDefinition:  1,
								MethodValue:       200,
								MeasurementSystem: 2,
								Reliability:       3,
							},
							{
								MethodDefinition:  3,
								MethodValue:       206,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  4,
								MethodValue:       216,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  5,
								MethodValue:       208,
								MeasurementSystem: 1,
								Reliability:       3,
							},
							{
								MethodDefinition:  6,
								MethodValue:       54,
								MeasurementSystem: 0,
								Reliability:       0,
							},
						},
					},
				},
			},
			bin: []byte{
				0x01,             // version
				0x00, 0x00, 0x00, // flags
				0x01, 0x00, 0x00, 0x00, 0x28, 0x32, 0x4C, 0x23, 0x05, 0x01, 0xC8, 0x23,
				0x03, 0xCE, 0x13, 0x04, 0xD8, 0x13, 0x05, 0xD0, 0x13, 0x06, 0x36, 0x00,
			},
			str: `Version=0 Flags=0x000000 LoudnessBases=[{EQSetID=0x0 DownmixID=0x0 DRCSetID=0x0 BsSamplePeakLevel=643 BsTruePeakLevel=588 MeasurementSystemForTP=0x2 ReliabilityForTP=0x3 MeasurementCount=0x5 Measurements=[{MethodDefinition=0x1 MethodValue=0xc8 MeasurementSystem=0x2 Reliability=0x3}, {MethodDefinition=0x3 MethodValue=0xce MeasurementSystem=0x1 Reliability=0x3}, {MethodDefinition=0x4 MethodValue=0xd8 MeasurementSystem=0x1 Reliability=0x3}, {MethodDefinition=0x5 MethodValue=0xd0 MeasurementSystem=0x1 Reliability=0x3}, {MethodDefinition=0x6 MethodValue=0x36 MeasurementSystem=0x0 Reliability=0x0}]}]`,
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
