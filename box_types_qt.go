package mp4

import "fmt"

/*************************** alac ****************************/

func BoxTypeAlac() BoxType {
	return StrToBoxType("alac")
}

func init() {
	AddBoxDef((*Alac)(nil))
	AddAnyTypeBoxDefEx(&AudioSampleEntry{}, BoxTypeAlac(), func(context Context) bool {
		return context.UnderStsd
	})
}

// Alac is ALACSpecificConfig entry
// https://github.com/macosforge/alac/blob/master/codec/ALACAudioTypes.h#L162
type Alac struct {
	FullBox `mp4:"0,extend"`

	FrameLength       uint32 `mp4:"1,size=32"`
	CompatibleVersion uint8  `mp4:"2,size=8"`
	BitDepth          uint8  `mp4:"3,size=8"`
	Pb                uint8  `mp4:"4,size=8"`
	Mb                uint8  `mp4:"5,size=8"`
	Kb                uint8  `mp4:"6,size=8"`
	NumChannels       uint8  `mp4:"7,size=8"`
	MaxRun            uint16 `mp4:"8,size=16"`
	MaxFrameByte      uint32 `mp4:"9,size=32"`
	AvgBitRate        uint32 `mp4:"10,size=32"`
	SampleRate        uint32 `mp4:"11,size=32"`
}

// GetType returns the BoxType
func (*Alac) GetType() BoxType {
	return BoxTypeAlac()
}

/*************************** chrm ****************************/

func BoxTypeChrm() BoxType {
	return StrToBoxType("chrm")
}

func init() {
	AddBoxDef((*Chrm)(nil))
}

// Chrm is AVC chrm box
type Chrm struct {
	Box
	X uint8 `mp4:"0,size=8"`
	Y uint8 `mp4:"1,size=8"`
}

// GetType returns the BoxType
func (*Chrm) GetType() BoxType {
	return BoxTypeChrm()
}

/*************************** ludt ****************************/

func BoxTypeLudt() BoxType {
	return StrToBoxType("ludt")
}
func BoxTypeTlou() BoxType {
	return StrToBoxType("tlou")
}
func BoxTypeAlou() BoxType {
	return StrToBoxType("alou")
}

func init() {
	AddBoxDef(&Ludt{})
	AddAnyTypeBoxDef(&LoudnessEntry{}, BoxTypeTlou())
	AddAnyTypeBoxDef(&LoudnessEntry{}, BoxTypeAlou())
}

// Ludt is Apple iTunes audio stream loudness box
type Ludt struct {
	Box
}

// GetType returns the BoxType
func (*Ludt) GetType() BoxType {
	return BoxTypeLudt()
}

type LoudnessEntry struct {
	AnyTypeBox        `mp4:"0,extend"`
	Version           uint8          `mp4:"0,size=8"`
	Flags             [3]byte        `mp4:"1,size=8"`
	LoudnessBaseCount uint8          `mp4:"2,size=8,nver=1,const=1"`
	LoudnessBases     []LoudnessBase `mp4:"3,len=dynamic"`
}

// GetFieldLength returns length of dynamic field
func (lou *LoudnessEntry) GetFieldLength(name string, ctx Context) uint {
	switch name {
	case "LoudnessBases":
		return uint(lou.LoudnessBaseCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=%s fieldName=%s", lou.Type.String(), name))
}

type LoudnessBase struct {
	BaseCustomFieldObject
	EQSetID                uint8                 `mp4:"0,size=8"`
	DownmixID              uint8                 `mp4:"1,size=10"`
	DRCSetID               uint8                 `mp4:"2,size=6"`
	BsSamplePeakLevel      uint16                `mp4:"3,size=12"`
	BsTruePeakLevel        uint16                `mp4:"4,size=12"`
	MeasurementSystemForTP uint8                 `mp4:"5,size=4"`
	ReliabilityForTP       uint8                 `mp4:"6,size=4"`
	MeasurementCount       uint8                 `mp4:"8,size=8"`
	Measurements           []LoudnessMeasurement `mp4:"9,size=24,len=dynamic"`
}

// GetFieldLength returns length of dynamic field
func (lou *LoudnessBase) GetFieldLength(name string, ctx Context) uint {
	switch name {
	case "Measurements":
		return uint(lou.MeasurementCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=loud_base fieldName=%s", name))
}

type LoudnessMeasurement struct {
	MethodDefinition  uint8 `mp4:"0,size=8"`
	MethodValue       uint8 `mp4:"1,size=8"`
	MeasurementSystem uint8 `mp4:"2,size=4"`
	Reliability       uint8 `mp4:"3,size=4"`
}

/*************************** swre ****************************/

func BoxTypeSwre() BoxType {
	return StrToBoxType("swre")
}

func init() {
	AddBoxDef(&Swre{})
}

// Swre is thef name and version number of the software that generated this movie
type Swre struct {
	FullBox `mp4:"0,extend"`
	// TODO: Meaning of these two bytes is still unknown.
	Unknown [2]byte `mp4:"1,size=8"`
	Version string  `mp4:"2"`
}

// GetType returns the BoxType
func (*Swre) GetType() BoxType {
	return BoxTypeSwre()
}
