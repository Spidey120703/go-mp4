package mp4

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type BoxPath []BoxType

func (lhs BoxPath) compareWith(rhs BoxPath) (forwardMatch bool, match bool) {
	if len(lhs) > len(rhs) {
		return false, false
	}
	for i := 0; i < len(lhs); i++ {
		if !lhs[i].MatchWith(rhs[i]) {
			return false, false
		}
	}
	if len(lhs) < len(rhs) {
		return true, false
	}
	return false, true
}

type ReadHandle struct {
	Params      []interface{}
	BoxInfo     BoxInfo
	Path        BoxPath
	ReadPayload func() (box IBox, n uint64, err error)
	ReadData    func(io.Writer) (n uint64, err error)
	Expand      func(params ...interface{}) (vals []interface{}, err error)
}

type ReadHandler func(handle *ReadHandle) (val interface{}, err error)

func ReadBoxStructure(r io.ReadSeeker, handler ReadHandler, params ...interface{}) ([]interface{}, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	return readBoxStructure(r, 0, true, nil, Context{}, handler, params)
}

// ReadBoxStructureWithContext reads and traverses the MP4 box structure
// using the provided parsing context.
//
// Unlike ReadBoxStructure, this function allows callers to supply an
// existing Context, which is useful for stateful parsing scenarios such as
// fragmented MP4, track-aware processing, or encryption-related metadata
// handling.
func ReadBoxStructureWithContext(r io.ReadSeeker, ctx Context, handler ReadHandler, params ...interface{}) ([]interface{}, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	return readBoxStructure(r, 0, true, nil, ctx, handler, params)
}

func ReadBoxStructureFromInternal(r io.ReadSeeker, bi *BoxInfo, handler ReadHandler, params ...interface{}) (interface{}, error) {
	return readBoxStructureFromInternal(r, bi, nil, handler, params)
}

func readBoxStructureFromInternal(r io.ReadSeeker, bi *BoxInfo, path BoxPath, handler ReadHandler, params []interface{}) (interface{}, error) {
	if _, err := bi.SeekToPayload(r); err != nil {
		return nil, err
	}

	// check comatible-brands
	if len(path) == 0 && bi.Type == BoxTypeFtyp() {
		var ftyp Ftyp
		if _, err := Unmarshal(r, bi.Size-bi.HeaderSize, &ftyp, bi.Context); err != nil {
			return nil, err
		}
		if ftyp.HasCompatibleBrand(BrandQT()) {
			bi.IsQuickTimeCompatible = true
		}
		if _, err := bi.SeekToPayload(r); err != nil {
			return nil, err
		}
	}

	// parse numbered ilst items after keys box by saving EntryCount field to context
	if bi.Type == BoxTypeKeys() {
		var keys Keys
		if _, err := Unmarshal(r, bi.Size-bi.HeaderSize, &keys, bi.Context); err != nil {
			return nil, err
		}
		bi.QuickTimeKeysMetaEntryCount = int(keys.EntryCount)
		if _, err := bi.SeekToPayload(r); err != nil {
			return nil, err
		}
	}

	ctx := bi.Context
	if bi.Type == BoxTypeWave() {
		ctx.UnderWave = true
	} else if bi.Type == BoxTypeIlst() {
		ctx.UnderIlst = true
	} else if bi.UnderIlst && !bi.UnderIlstMeta && IsIlstMetaBoxType(bi.Type) {
		ctx.UnderIlstMeta = true
		if bi.Type == StrToBoxType("----") {
			ctx.UnderIlstFreeMeta = true
		}
	} else if bi.Type == BoxTypeUdta() {
		ctx.UnderUdta = true
	} else if bi.Type == BoxTypeStsd() {
		ctx.UnderStsd = true
	} else if ctx.UnderStsd {
		// handle box type collision between stsd sample entry and its nested boxes
		ctx.UnderStsd = false
	}

	if bi.Type == BoxTypeTrak() {
		// read tkhd box size
		var tfhdSize uint32
		_ = binary.Read(r, binary.BigEndian, &tfhdSize)
		_, _ = r.Seek(4, io.SeekCurrent)

		// read track id from tkhd
		var tkhd Tkhd
		if _, err := Unmarshal(r, uint64(tfhdSize), &tkhd, bi.Context); err != nil {
			return nil, err
		}
		ctx.TrackID = tkhd.TrackID
		if _, err := bi.SeekToPayload(r); err != nil {
			return nil, err
		}
	} else if bi.Type == BoxTypeTraf() {
		// read tfhd box size
		var tfhdSize uint32
		_ = binary.Read(r, binary.BigEndian, &tfhdSize)
		_, _ = r.Seek(4, io.SeekCurrent)

		// read track id from tfhd
		var tfhd Tfhd
		if _, err := Unmarshal(r, uint64(tfhdSize), &tfhd, bi.Context); err != nil {
			return nil, err
		}
		ctx.TrackID = tfhd.TrackID
		if _, err := bi.SeekToPayload(r); err != nil {
			return nil, err
		}
	} else if bi.Type == BoxTypeTenc() {
		// register the current track ID with its corresponding tenc box
		var tenc Tenc
		if _, err := Unmarshal(r, bi.Size-bi.HeaderSize, &tenc, bi.Context); err != nil {
			return nil, err
		}
		if ctx.Crypto != nil && ctx.Crypto.TencRegistry != nil {
			ctx.Crypto.TencRegistry[ctx.TrackID] = &tenc
		}
		if _, err := bi.SeekToPayload(r); err != nil {
			return nil, err
		}
	}

	newPath := make(BoxPath, len(path)+1)
	copy(newPath, path)
	newPath[len(path)] = bi.Type

	h := &ReadHandle{
		Params:  params,
		BoxInfo: *bi,
		Path:    newPath,
	}

	var childrenOffset uint64

	h.ReadPayload = func() (IBox, uint64, error) {
		if _, err := bi.SeekToPayload(r); err != nil {
			return nil, 0, err
		}

		box, n, err := UnmarshalAny(r, bi.Type, bi.Size-bi.HeaderSize, bi.Context)
		if err != nil {
			return nil, 0, err
		}
		childrenOffset = bi.Offset + bi.HeaderSize + n
		return box, n, nil
	}

	h.ReadData = func(w io.Writer) (uint64, error) {
		if _, err := bi.SeekToPayload(r); err != nil {
			return 0, err
		}

		size := bi.Size - bi.HeaderSize
		if _, err := io.CopyN(w, r, int64(size)); err != nil {
			return 0, err
		}
		return size, nil
	}

	h.Expand = func(params ...interface{}) ([]interface{}, error) {
		if childrenOffset == 0 {
			if _, err := bi.SeekToPayload(r); err != nil {
				return nil, err
			}

			_, n, err := UnmarshalAny(r, bi.Type, bi.Size-bi.HeaderSize, bi.Context)
			if err != nil {
				return nil, err
			}
			childrenOffset = bi.Offset + bi.HeaderSize + n
		} else {
			if _, err := r.Seek(int64(childrenOffset), io.SeekStart); err != nil {
				return nil, err
			}
		}

		childrenSize := bi.Offset + bi.Size - childrenOffset
		return readBoxStructure(r, childrenSize, false, newPath, ctx, handler, params)
	}

	if val, err := handler(h); err != nil {
		return nil, err
	} else if _, err := bi.SeekToEnd(r); err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func readBoxStructure(r io.ReadSeeker, totalSize uint64, isRoot bool, path BoxPath, ctx Context, handler ReadHandler, params []interface{}) ([]interface{}, error) {
	vals := make([]interface{}, 0, 8)

	for isRoot || totalSize >= SmallHeaderSize {
		bi, err := ReadBoxInfo(r)
		if isRoot && err == io.EOF {
			return vals, nil
		} else if err != nil {
			return nil, err
		}

		if !isRoot && bi.Size > totalSize {
			return nil, fmt.Errorf("too large box size: type=%s, size=%d, actualBufSize=%d", bi.Type.String(), bi.Size, totalSize)
		}
		totalSize -= bi.Size

		bi.Context = ctx

		val, err := readBoxStructureFromInternal(r, bi, path, handler, params)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)

		if bi.IsQuickTimeCompatible {
			ctx.IsQuickTimeCompatible = true
		}

		// preserve keys entry count on context for subsequent ilst number item box
		if bi.Type == BoxTypeKeys() {
			ctx.QuickTimeKeysMetaEntryCount = bi.QuickTimeKeysMetaEntryCount
		}
	}

	if totalSize != 0 && !ctx.IsQuickTimeCompatible {
		return nil, errors.New("unexpected EOF")
	}

	return vals, nil
}
