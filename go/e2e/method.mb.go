package e2e

import "github.com/orbs-network/membuffers/go"

/*
message Method {
	string name = 1;
	repeated MethodCallArgument arg = 2;
}
*/

type Method struct {
	_Message membuffers.Message
}

var m_Method_Scheme = []membuffers.FieldType{membuffers.TypeString,membuffers.TypeMessageArray}
var m_Method_Unions = [][]membuffers.FieldType{{}}

func ReadMethod(buf []byte) *Method {
	x := &Method{}
	x._Message.Init(buf, membuffers.Offset(len(buf)), m_Method_Scheme, m_Method_Unions)
	return x
}

func (x *Method) _RawBuffer() []byte {
	return x._Message.RawBuffer()
}

func (x *Method) Name() string {
	return x._Message.GetString(0)
}

func (x *Method) _RawBuffer_Name() []byte {
	return x._Message.RawBufferForField(0, 0)
}

func (x *Method) ArgIterator() *Method_ArgIterator {
	return &Method_ArgIterator{_Iterator: x._Message.GetMessageArrayIterator(1)}
}

type Method_ArgIterator struct {
	_Iterator *membuffers.Iterator
}

func (i *Method_ArgIterator) HasNext() bool {
	return i._Iterator.HasNext()
}

func (i *Method_ArgIterator) NextArg() *MethodCallArgument {
	b, s := i._Iterator.NextMessage()
	return ReadMethodCallArgument(b[:s])
}

func (x *Method) _RawBuffer_ArgArray() []byte {
	return x._Message.RawBufferForField(1, 0)
}

/*
message MethodCallArgument {
	oneof type {
		uint32 num = 1;
		string str = 2;
		bytes data = 3;
	}
}
*/

type MethodCallArgument struct {
	_Message membuffers.Message
}

var m_MethodCallArgument_Scheme = []membuffers.FieldType{membuffers.TypeUnion}
var m_MethodCallArgument_Unions = [][]membuffers.FieldType{{membuffers.TypeUint32,membuffers.TypeString,membuffers.TypeBytes}}

func ReadMethodCallArgument(buf []byte) *MethodCallArgument {
	x := &MethodCallArgument{}
	x._Message.Init(buf, membuffers.Offset(len(buf)), m_MethodCallArgument_Scheme, m_MethodCallArgument_Unions)
	return x
}

func (x *MethodCallArgument) _RawBuffer() []byte {
	return x._Message.RawBuffer()
}

type MethodCallArgument_Type uint16

const (
	MethodCallArgument_Type_Num MethodCallArgument_Type = 0
	MethodCallArgument_Type_Str MethodCallArgument_Type = 1
	MethodCallArgument_Type_Data MethodCallArgument_Type = 2
)

func (x *MethodCallArgument) Type() MethodCallArgument_Type {
	return MethodCallArgument_Type(x._Message.GetUint16(0))
}

func (x *MethodCallArgument) IsType_Num() bool {
	is, _ := x._Message.IsUnionIndex(0, 0, 0)
	return is
}

func (x *MethodCallArgument) Type_Num() uint32 {
	_, off := x._Message.IsUnionIndex(0, 0, 0)
	return x._Message.GetUint32InOffset(off)
}

func (x *MethodCallArgument) IsType_Str() bool {
	is, _ := x._Message.IsUnionIndex(0, 0, 1)
	return is
}

func (x *MethodCallArgument) Type_Str() string {
	_, off := x._Message.IsUnionIndex(0, 0, 1)
	return x._Message.GetStringInOffset(off)
}

func (x *MethodCallArgument) IsType_Data() bool {
	is, _ := x._Message.IsUnionIndex(0, 0, 2)
	return is
}

func (x *MethodCallArgument) Type_Data() []byte {
	_, off := x._Message.IsUnionIndex(0, 0, 2)
	return x._Message.GetBytesInOffset(off)
}

func (x *MethodCallArgument) _RawBuffer_Type() []byte {
	return x._Message.RawBufferForField(0, 0)
}