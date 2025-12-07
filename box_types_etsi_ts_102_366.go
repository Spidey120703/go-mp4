package mp4

/*************************** ac-3 ****************************/

// https://www.etsi.org/deliver/etsi_ts/102300_102399/102366/01.04.01_60/ts_102366v010401p.pdf

func BoxTypeAC3() BoxType { return StrToBoxType("ac-3") }

func init() {
	AddAnyTypeBoxDef(&AudioSampleEntry{}, BoxTypeAC3())
}

/*************************** dac3 ****************************/

// https://www.etsi.org/deliver/etsi_ts/102300_102399/102366/01.04.01_60/ts_102366v010401p.pdf

func BoxTypeDAC3() BoxType { return StrToBoxType("dac3") }

func init() {
	AddBoxDef(&Dac3{})
}

type Dac3 struct {
	Box
	Fscod       uint8 `mp4:"0,size=2"`
	Bsid        uint8 `mp4:"1,size=5"`
	Bsmod       uint8 `mp4:"2,size=3"`
	Acmod       uint8 `mp4:"3,size=3"`
	LfeOn       uint8 `mp4:"4,size=1"`
	BitRateCode uint8 `mp4:"5,size=5"`
	Reserved    uint8 `mp4:"6,size=5,const=0"`
}

func (Dac3) GetType() BoxType {
	return BoxTypeDAC3()
}

/*************************** ec-3 ****************************/

// https://www.etsi.org/deliver/etsi_ts/103400_103499/103420/01.02.01_60/ts_103420v010201p.pdf

func BoxTypeEC3() BoxType { return StrToBoxType("ec-3") }

func init() {
	AddAnyTypeBoxDef(&AudioSampleEntry{}, BoxTypeEC3())
}

/*************************** dec3 ****************************/

func BoxTypeDEC3() BoxType { return StrToBoxType("dec3") }

func init() {
	AddBoxDef(&Dec3{})
}

type Dec3 struct {
	Box
	DataRate              uint16   `mp4:"0,size=13"`
	NumIndSub             uint8    `mp4:"1,size=3"`
	IndSub                []IndSub `mp4:"2,len=dynamic"`
	Reserved1             uint8    `mp4:"3,size=7"`
	FlagEC3ExtensionTypeA uint8    `mp4:"4,size=1"`
	ComplexityIndexTypeA  uint8    `mp4:"5,size=8"`
	Reserved2             []byte   `mp4:"6,size=8"`
}

// GetFieldLength returns length of dynamic field
func (dec3 Dec3) GetFieldLength(name string, ctx Context) uint {
	switch name {
	case "IndSub":
		return uint(dec3.NumIndSub) + 1
	}
	return 0
}

type IndSub struct {
	BaseCustomFieldObject
	Fscod     uint8  `mp4:"0,size=2"`
	Bsid      uint8  `mp4:"1,size=5"`
	Reserved1 uint8  `mp4:"2,size=1,const=0"`
	Asvc      uint8  `mp4:"3,size=1"`
	Bsmod     uint8  `mp4:"4,size=3"`
	Acmod     uint8  `mp4:"5,size=3"`
	LfeOn     uint8  `mp4:"6,size=1"`
	Reserved2 uint8  `mp4:"7,size=3,const=0"`
	NumDepSub uint8  `mp4:"8,size=4"`
	ChanLoc   uint16 `mp4:"9,size=9,opt=dynamic"`
	Reserved3 uint8  `mp4:"10,size=1,const=0,opt=dynamic"`
}

func (is IndSub) IsOptFieldEnabled(name string, ctx Context) bool {
	switch name {
	case "ChanLoc":
		return is.NumDepSub > 0
	case "Reserved3":
		return is.ChanLoc <= 0
	}
	return false
}

func (Dec3) GetType() BoxType {
	return BoxTypeDEC3()
}
