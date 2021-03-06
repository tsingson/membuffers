// Copyright 2018 the membuffers authors
// This file is part of the membuffers library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package membuffers

import "math/big"

type MessageWriter interface {
	HexDump(prefix string, offsetFromStart Offset) (err error)
	Write(buf []byte) (err error)
	GetSize() Offset
	CalcRequiredSize() Offset
}

type InternalBuilder struct {
	size     Offset
	building bool
}

func (w *InternalBuilder) Reset() {
	w.size = 0
}

// override me
func (w *InternalBuilder) Write(buf []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &ErrBufferOverrun{}
		}
	}()
	w.Reset()
	return nil
}

// instead of implementing Write with the actual fields, this allows to take a ready-made raw buffer
func (w *InternalBuilder) WriteOverrideWithRawBuffer(buf []byte, overrideWithRawBuffer []byte) error {
	if buf != nil {
		if len(buf) < len(overrideWithRawBuffer) {
			return &ErrBufferOverrun{}
		}
		copy(buf, overrideWithRawBuffer)
	}
	w.size = Offset(len(overrideWithRawBuffer))
	return nil
}

func (w *InternalBuilder) CalcRequiredSize() Offset {
	w.Write(nil)
	return w.GetSize()
}

func (w *InternalBuilder) GetSize() Offset {
	return w.size
}

func (w *InternalBuilder) NotifyBuildStart() {
	if w.building {
		panic("concurrent build: builder is not thread safe (started building from multiple goroutines)")
	}
	w.building = true // not a mutex for performance reasons
}

func (w *InternalBuilder) NotifyBuildEnd() {
	w.building = false
}

func (w *InternalBuilder) WriteBool(buf []byte, v bool) {
	w.size = alignOffsetToType(w.size, TypeBool)
	if buf != nil {
		WriteBool(buf[w.size:], v)
	}
	w.size += FieldSizes[TypeBool]
}

func (w *InternalBuilder) WriteUint8(buf []byte, v uint8) {
	w.size = alignOffsetToType(w.size, TypeUint8)
	if buf != nil {
		WriteUint8(buf[w.size:], v)
	}
	w.size += FieldSizes[TypeUint8]
}

func (w *InternalBuilder) WriteUint16(buf []byte, v uint16) {
	w.size = alignOffsetToType(w.size, TypeUint16)
	if buf != nil {
		WriteUint16(buf[w.size:], v)
	}
	w.size += FieldSizes[TypeUint16]
}

func (w *InternalBuilder) WriteUint32(buf []byte, v uint32) {
	w.size = alignOffsetToType(w.size, TypeUint32)
	if buf != nil {
		WriteUint32(buf[w.size:], v)
	}
	w.size += FieldSizes[TypeUint32]
}

func (w *InternalBuilder) WriteUint64(buf []byte, v uint64) {
	w.size = alignOffsetToType(w.size, TypeUint64)
	if buf != nil {
		WriteUint64(buf[w.size:], v)
	}
	w.size += FieldSizes[TypeUint64]
}

func (w *InternalBuilder) WriteBytes(buf []byte, v []byte) {
	w.size = alignOffsetToType(w.size, TypeBytes)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v)))
	}
	w.size += FieldSizes[TypeBytes]
	w.size = alignDynamicFieldContentOffset(w.size, TypeBytes)
	if v != nil {
		if buf != nil {
			copy(buf[w.size:], v)
		}
		w.size += Offset(len(v))
	}
}

func (w *InternalBuilder) WriteBytes20(buf []byte, v [20]byte) {
	w.size = alignOffsetToType(w.size, TypeBytes20)
	if buf != nil {
		WriteBytes20(buf[w.size:], v)
	}
	w.size += FieldSizes[TypeBytes20]
}

func (w *InternalBuilder) WriteBytes32(buf []byte, v [32]byte) {
	w.size = alignOffsetToType(w.size, TypeBytes32)
	if buf != nil {
		WriteBytes32(buf[w.size:], v)
	}
	w.size += FieldSizes[TypeBytes32]
}

func (w *InternalBuilder) WriteUint256(buf []byte, v *big.Int) {
	w.size = alignOffsetToType(w.size, TypeUint256)
	if buf != nil {
		WriteUint256(buf[w.size:], v)
	}
	w.size += FieldSizes[TypeUint256]
}

func (w *InternalBuilder) WriteString(buf []byte, v string) {
	w.WriteBytes(buf, []byte(v))
}

func (w *InternalBuilder) WriteUnionIndex(buf []byte, unionIndex uint16) {
	w.size = alignOffsetToType(w.size, TypeUnion)
	if buf != nil {
		WriteUnionType(buf[w.size:], unionIndex)
	}
	w.size += FieldSizes[TypeUnion]
}

func (w *InternalBuilder) WriteBoolArray(buf []byte, v []bool) {
	w.size = alignOffsetToType(w.size, TypeBoolArray)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v))*FieldSizes[TypeBool])
	}
	w.size += FieldSizes[TypeBoolArray]
	w.size = alignDynamicFieldContentOffset(w.size, TypeBoolArray)
	for _, vv := range v {
		w.WriteBool(buf, vv)
	}
}

func (w *InternalBuilder) WriteUint8Array(buf []byte, v []uint8) {
	w.size = alignOffsetToType(w.size, TypeUint8Array)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v))*FieldSizes[TypeUint8])
	}
	w.size += FieldSizes[TypeUint8Array]
	w.size = alignDynamicFieldContentOffset(w.size, TypeUint8Array)
	for _, vv := range v {
		w.WriteUint8(buf, vv)
	}
}

func (w *InternalBuilder) WriteUint16Array(buf []byte, v []uint16) {
	w.size = alignOffsetToType(w.size, TypeUint16Array)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v))*FieldSizes[TypeUint16])
	}
	w.size += FieldSizes[TypeUint16Array]
	w.size = alignDynamicFieldContentOffset(w.size, TypeUint16Array)
	for _, vv := range v {
		w.WriteUint16(buf, vv)
	}
}

func (w *InternalBuilder) WriteUint32Array(buf []byte, v []uint32) {
	w.size = alignOffsetToType(w.size, TypeUint32Array)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v))*FieldSizes[TypeUint32])
	}
	w.size += FieldSizes[TypeUint32Array]
	w.size = alignDynamicFieldContentOffset(w.size, TypeUint32Array)
	for _, vv := range v {
		w.WriteUint32(buf, vv)
	}
}

func (w *InternalBuilder) WriteUint64Array(buf []byte, v []uint64) {
	w.size = alignOffsetToType(w.size, TypeUint64Array)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v))*FieldSizes[TypeUint64])
	}
	w.size += FieldSizes[TypeUint64Array]
	w.size = alignDynamicFieldContentOffset(w.size, TypeUint64Array)
	for _, vv := range v {
		w.WriteUint64(buf, vv)
	}
}

func (w *InternalBuilder) WriteBytesArray(buf []byte, v [][]byte) {
	w.size = alignOffsetToType(w.size, TypeBytesArray)
	sizePlaceholderOffset := w.size
	w.size += FieldSizes[TypeBytesArray]
	w.size = alignDynamicFieldContentOffset(w.size, TypeBytesArray)
	contentSizeStartOffset := w.size
	for _, vv := range v {
		w.WriteBytes(buf, vv)
	}
	contentSize := w.size - contentSizeStartOffset
	if buf != nil {
		WriteOffset(buf[sizePlaceholderOffset:], contentSize)
	}
}

func (w *InternalBuilder) WriteBytes20Array(buf []byte, v [][20]byte) {
	w.size = alignOffsetToType(w.size, TypeBytes20Array)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v)*20))
	}
	w.size += FieldSizes[TypeBytes20Array]
	for _, vv := range v {
		w.WriteBytes20(buf, vv)
	}
}

func (w *InternalBuilder) WriteBytes32Array(buf []byte, v [][32]byte) {
	w.size = alignOffsetToType(w.size, TypeBytes32Array)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v)*32))
	}
	w.size += FieldSizes[TypeBytes32Array]
	for _, vv := range v {
		w.WriteBytes32(buf, vv)
	}
}

func (w *InternalBuilder) WriteUint256Array(buf []byte, v []*big.Int) {
	w.size = alignOffsetToType(w.size, TypeUint256Array)
	if buf != nil {
		WriteOffset(buf[w.size:], Offset(len(v)*32))
	}
	w.size += FieldSizes[TypeUint256Array]
	for _, vv := range v {
		w.WriteUint256(buf, vv)
	}
}

func (w *InternalBuilder) WriteStringArray(buf []byte, v []string) {
	w.size = alignOffsetToType(w.size, TypeStringArray)
	sizePlaceholderOffset := w.size
	w.size += FieldSizes[TypeStringArray]
	w.size = alignDynamicFieldContentOffset(w.size, TypeStringArray)
	contentSizeStartOffset := w.size
	for _, vv := range v {
		w.WriteString(buf, vv)
	}
	contentSize := w.size - contentSizeStartOffset
	if buf != nil {
		WriteOffset(buf[sizePlaceholderOffset:], contentSize)
	}
}

func (w *InternalBuilder) WriteMessage(buf []byte, v MessageWriter) (err error) {
	w.size = alignOffsetToType(w.size, TypeMessage)
	sizePlaceholderOffset := w.size
	w.size += FieldSizes[TypeMessage]
	w.size = alignDynamicFieldContentOffset(w.size, TypeMessage)
	if buf != nil {
		err = v.Write(buf[w.size:])
	} else {
		err = v.Write(nil)
	}
	contentSize := v.GetSize()
	w.size += contentSize
	if buf != nil {
		WriteOffset(buf[sizePlaceholderOffset:], contentSize)
	}
	return
}

func (w *InternalBuilder) WriteMessageArray(buf []byte, v []MessageWriter) (err error) {
	w.size = alignOffsetToType(w.size, TypeMessageArray)
	sizePlaceholderOffset := w.size
	w.size += FieldSizes[TypeMessageArray]
	w.size = alignDynamicFieldContentOffset(w.size, TypeMessageArray)
	contentSizeStartOffset := w.size
	for _, vv := range v {
		err = w.WriteMessage(buf, vv)
		if err != nil {
			return
		}
	}
	contentSize := w.size - contentSizeStartOffset
	if buf != nil {
		WriteOffset(buf[sizePlaceholderOffset:], contentSize)
	}
	return nil
}
