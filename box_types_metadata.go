package mp4

import (
	"fmt"

	"github.com/Spidey120703/go-mp4/internal/util"
)

/*************************** hdlr ****************************/

func init() {
	AddBoxDefEx(&MetadataHandlerBox{}, isUnderUdta)
}

type MetadataHandlerBox struct {
	FullBox       `mp4:"0,extend"`
	ComponentType uint32   `mp4:"1,size=32"`
	HandlerType   [4]byte  `mp4:"2,size=8,string"`
	Name          [14]byte `mp4:"3,size=8,string"`
}

// GetType returns the BoxType
func (*MetadataHandlerBox) GetType() BoxType {
	return BoxTypeHdlr()
}

/*************************** ilst ****************************/

func BoxTypeIlst() BoxType { return StrToBoxType("ilst") }
func BoxTypeData() BoxType { return StrToBoxType("data") }

var ilstMetaBoxTypes = []BoxType{
	StrToBoxType("----"),
	StrToBoxType("@PST"),
	StrToBoxType("@ppi"),
	StrToBoxType("@pti"),
	StrToBoxType("@sti"),
	StrToBoxType("AACR"),
	StrToBoxType("CDEK"),
	StrToBoxType("CDET"),
	StrToBoxType("GUID"),
	StrToBoxType("VERS"),
	StrToBoxType("aART"),
	StrToBoxType("akID"),
	StrToBoxType("albm"),
	StrToBoxType("apID"),
	StrToBoxType("atID"),
	StrToBoxType("auth"),
	StrToBoxType("catg"),
	StrToBoxType("cmID"),
	StrToBoxType("cnID"),
	StrToBoxType("covr"),
	StrToBoxType("cpil"),
	StrToBoxType("cprt"),
	StrToBoxType("desc"),
	StrToBoxType("disk"),
	StrToBoxType("dscp"),
	StrToBoxType("egid"),
	StrToBoxType("geID"),
	StrToBoxType("gnre"),
	StrToBoxType("grup"),
	StrToBoxType("gshh"),
	StrToBoxType("gspm"),
	StrToBoxType("gspu"),
	StrToBoxType("gssd"),
	StrToBoxType("gsst"),
	StrToBoxType("gstd"),
	StrToBoxType("hdvd"),
	StrToBoxType("itnu"),
	StrToBoxType("keyw"),
	StrToBoxType("ldes"),
	StrToBoxType("ownr"),
	StrToBoxType("pcst"),
	StrToBoxType("perf"),
	StrToBoxType("pgap"),
	StrToBoxType("plID"),
	StrToBoxType("prID"),
	StrToBoxType("purd"),
	StrToBoxType("purl"),
	StrToBoxType("rate"),
	StrToBoxType("rldt"),
	StrToBoxType("rtng"),
	StrToBoxType("sdes"),
	StrToBoxType("sfID"),
	StrToBoxType("shwm"),
	StrToBoxType("snal"),
	StrToBoxType("soaa"),
	StrToBoxType("soal"),
	StrToBoxType("soar"),
	StrToBoxType("soco"),
	StrToBoxType("sonm"),
	StrToBoxType("sosn"),
	StrToBoxType("stik"),
	StrToBoxType("titl"),
	StrToBoxType("tmpo"),
	StrToBoxType("tnal"),
	StrToBoxType("trkn"),
	StrToBoxType("tven"),
	StrToBoxType("tves"),
	StrToBoxType("tvnn"),
	StrToBoxType("tvsh"),
	StrToBoxType("tvsn"),
	StrToBoxType("xid "),
	StrToBoxType("yrrc"),
	StrToBoxType("data"),
	{0xA9, 'A', 'R', 'T'},
	{0xA9, 'a', 'l', 'b'},
	{0xA9, 'a', 'r', 'd'},
	{0xA9, 'a', 'r', 'g'},
	{0xA9, 'a', 'u', 't'},
	{0xA9, 'c', 'm', 't'},
	{0xA9, 'c', 'o', 'm'},
	{0xA9, 'c', 'o', 'n'},
	{0xA9, 'c', 'p', 'y'},
	{0xA9, 'd', 'a', 'y'},
	{0xA9, 'd', 'e', 's'},
	{0xA9, 'd', 'i', 'r'},
	{0xA9, 'e', 'n', 'c'},
	{0xA9, 'g', 'e', 'n'},
	{0xA9, 'g', 'r', 'p'},
	{0xA9, 'l', 'y', 'r'},
	{0xA9, 'm', 'v', 'c'},
	{0xA9, 'm', 'v', 'i'},
	{0xA9, 'm', 'v', 'n'},
	{0xA9, 'n', 'a', 'm'},
	{0xA9, 'n', 'r', 't'},
	{0xA9, 'o', 'p', 'e'},
	{0xA9, 'p', 'r', 'd'},
	{0xA9, 'p', 'u', 'b'},
	{0xA9, 's', 'n', 'e'},
	{0xA9, 's', 'o', 'l'},
	{0xA9, 's', 't', '3'},
	{0xA9, 't', 'o', 'o'},
	{0xA9, 't', 'r', 'k'},
	{0xA9, 'w', 'r', 'k'},
	{0xA9, 'w', 'r', 't'},
	{0xA9, 'x', 'p', 'd'},
	{0xA9, 'x', 'y', 'z'},
}

func IsIlstMetaBoxType(boxType BoxType) bool {
	for _, bt := range ilstMetaBoxTypes {
		if boxType == bt {
			return true
		}
	}
	return false
}

func init() {
	AddBoxDef(&Ilst{})
	AddBoxDefEx(&Data{}, isUnderIlstMeta)
	for _, bt := range ilstMetaBoxTypes {
		AddAnyTypeBoxDefEx(&IlstMetaContainer{}, bt, isIlstMetaContainer)
	}
	AddAnyTypeBoxDefEx(&StringData{}, StrToBoxType("mean"), isUnderIlstFreeFormat)
	AddAnyTypeBoxDefEx(&StringData{}, StrToBoxType("name"), isUnderIlstFreeFormat)
}

type Ilst struct {
	Box
}

// GetType returns the BoxType
func (*Ilst) GetType() BoxType {
	return BoxTypeIlst()
}

type IlstMetaContainer struct {
	AnyTypeBox
}

func isIlstMetaContainer(ctx Context) bool {
	return ctx.UnderIlst && !ctx.UnderIlstMeta
}

const (
	DatatypeReserved           = 0
	DataTypeUTF8               = 1
	DataTypeUTF16              = 2
	DataTypeSJIS               = 3
	DataTypeUTF8Sort           = 4
	DataTypeUTF16Sort          = 5
	DataTypeJPEG               = 13
	DataTypePNG                = 14
	DataTypeInt                = 21
	DataTypeUint               = 22
	DataTypeFloat32            = 23
	DataTypeFloat64            = 24
	DataTypeBMP                = 27
	DataTypeQTMetadataAtom     = 28
	DataTypeInt8               = 65
	DataTypeInt16              = 66
	DataTypeInt32              = 67
	DataTypePointF32           = 70
	DataTypeDimensionsF32      = 71
	DataTypeRectF32            = 72
	DataTypeInt64              = 74
	DataTypeUint8              = 75
	DataTypeUint16             = 76
	DataTypeUint32             = 77
	DataTypeUint64             = 78
	DataTypeAffineTransformF64 = 79
)

// Data is a Value BoxType
// https://developer.apple.com/documentation/quicktime-file-format/value_atom
type Data struct {
	Box
	DataType uint32 `mp4:"0,size=32"`
	DataLang uint32 `mp4:"1,size=32"`
	Data     []byte `mp4:"2,size=8"`
}

// GetType returns the BoxType
func (*Data) GetType() BoxType {
	return BoxTypeData()
}

func isUnderIlstMeta(ctx Context) bool {
	return ctx.UnderIlstMeta
}

// StringifyField returns field value as string
func (data *Data) StringifyField(name string, indent string, depth int, ctx Context) (string, bool) {
	switch name {
	case "DataType":
		switch data.DataType {
		case DatatypeReserved:
			return "BINARY", true
		case DataTypeUTF8, DataTypeUTF8Sort:
			return "UTF8", true
		case DataTypeUTF16, DataTypeUTF16Sort:
			return "UTF16", true
		case DataTypeSJIS:
			return "SJIS", true
		case DataTypeJPEG:
			return "JPEG", true
		case DataTypePNG:
			return "PNG", true
		case DataTypeBMP:
			return "BMP", true
		case DataTypeInt:
			return "INT", true
		case DataTypeUint:
			return "UINT", true
		case DataTypeInt8:
			return "INT8", true
		case DataTypeUint8:
			return "UINT8", true
		case DataTypeInt16:
			return "INT16", true
		case DataTypeUint16:
			return "UINT16", true
		case DataTypeInt32:
			return "INT32", true
		case DataTypeUint32:
			return "UINT32", true
		case DataTypeInt64:
			return "INT64", true
		case DataTypeUint64:
			return "UINT64", true
		case DataTypeFloat32:
			return "FLOAT32", true
		case DataTypeFloat64:
			return "FLOAT64", true
		}
	case "Data":
		switch data.DataType {
		case DataTypeUTF8:
			return fmt.Sprintf("\"%s\"", util.EscapeUnprintables(string(data.Data))), true
		}
	}
	return "", false
}

type StringData struct {
	AnyTypeBox
	Data []byte `mp4:"0,size=8"`
}

// StringifyField returns field value as string
func (sd *StringData) StringifyField(name string, indent string, depth int, ctx Context) (string, bool) {
	if name == "Data" {
		return fmt.Sprintf("\"%s\"", util.EscapeUnprintables(string(sd.Data))), true
	}
	return "", false
}

/*************************** numbered items ****************************/

// Item is a numbered item under an item list atom
// https://developer.apple.com/documentation/quicktime-file-format/metadata_item_list_atom/item_list
type Item struct {
	AnyTypeBox
	Version  uint8   `mp4:"0,size=8"`
	Flags    [3]byte `mp4:"1,size=8"`
	ItemName []byte  `mp4:"2,size=8,len=4"`
	Data     Data    `mp4:"3"`
}

// StringifyField returns field value as string
func (i *Item) StringifyField(name string, indent string, depth int, ctx Context) (string, bool) {
	switch name {
	case "ItemName":
		return fmt.Sprintf("\"%s\"", util.EscapeUnprintables(string(i.ItemName))), true
	}
	return "", false
}

func isUnderIlstFreeFormat(ctx Context) bool {
	return ctx.UnderIlstFreeMeta
}

func BoxTypeKeys() BoxType { return StrToBoxType("keys") }

func init() {
	AddBoxDef(&Keys{})
}

/*************************** keys ****************************/

// Keys is the Keys BoxType
// https://developer.apple.com/documentation/quicktime-file-format/metadata_item_keys_atom
type Keys struct {
	FullBox    `mp4:"0,extend"`
	EntryCount int32 `mp4:"1,size=32"`
	Entries    []Key `mp4:"2,len=dynamic"`
}

// GetType implements the IBox interface and returns the BoxType
func (*Keys) GetType() BoxType {
	return BoxTypeKeys()
}

// GetFieldLength implements the ICustomFieldObject interface and returns the length of dynamic fields
func (k *Keys) GetFieldLength(name string, ctx Context) uint {
	switch name {
	case "Entries":
		return uint(k.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=keys fieldName=%s", name))
}

/*************************** key ****************************/

// Key is a key value field in the Keys BoxType
// https://developer.apple.com/documentation/quicktime-file-format/metadata_item_keys_atom/key_value_key_size-8
type Key struct {
	BaseCustomFieldObject
	KeySize      int32  `mp4:"0,size=32"`
	KeyNamespace []byte `mp4:"1,size=8,len=4"`
	KeyValue     []byte `mp4:"2,size=8,len=dynamic"`
}

// GetFieldLength implements the ICustomFieldObject interface and returns the length of dynamic fields
func (k *Key) GetFieldLength(name string, ctx Context) uint {
	switch name {
	case "KeyValue":
		// sizeOf(KeySize)+sizeOf(KeyNamespace) = 8 bytes
		return uint(k.KeySize) - 8
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=key fieldName=%s", name))
}

// StringifyField returns field value as string
func (k *Key) StringifyField(name string, indent string, depth int, ctx Context) (string, bool) {
	switch name {
	case "KeyNamespace":
		return fmt.Sprintf("\"%s\"", util.EscapeUnprintables(string(k.KeyNamespace))), true
	case "KeyValue":
		return fmt.Sprintf("\"%s\"", util.EscapeUnprintables(string(k.KeyValue))), true
	}
	return "", false
}
