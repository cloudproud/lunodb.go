package types

import (
	lunopb "github.com/cloudproud/lunodb.api/proto/types"
)

const (
	Any       = lunopb.Any
	Bool      = lunopb.Bool
	String    = lunopb.String
	Int8      = lunopb.Int8
	Int16     = lunopb.Int16
	Int32     = lunopb.Int32
	Int64     = lunopb.Int64
	Uint8     = lunopb.Uint8
	Uint16    = lunopb.Uint16
	Uint32    = lunopb.Uint32
	Uint64    = lunopb.Uint64
	Float32   = lunopb.Float32
	Float64   = lunopb.Float64
	Array     = lunopb.Array
	Object    = lunopb.Object
	Tuple     = lunopb.Tuple
	Time      = lunopb.Time
	Date      = lunopb.Date
	Bytes     = lunopb.Bytes
	UUID      = lunopb.UUID
	Timestamp = lunopb.Timestamp
	Inet      = lunopb.Inet
	Duration  = lunopb.Duration
	Record    = lunopb.Record
)

var BasicAny = &lunopb.Type{Kind: Any}
var BasicBool = &lunopb.Type{Kind: Bool}
var BasicString = &lunopb.Type{Kind: String}
var BasicInt8 = &lunopb.Type{Kind: Int8}
var BasicInt16 = &lunopb.Type{Kind: Int16}
var BasicInt32 = &lunopb.Type{Kind: Int32}
var BasicInt64 = &lunopb.Type{Kind: Int64}
var BasicUint8 = &lunopb.Type{Kind: Uint8}
var BasicUint16 = &lunopb.Type{Kind: Uint16}
var BasicUint32 = &lunopb.Type{Kind: Uint32}
var BasicUint64 = &lunopb.Type{Kind: Uint64}
var BasicFloat32 = &lunopb.Type{Kind: Float32}
var BasicFloat64 = &lunopb.Type{Kind: Float64}
var BasicObject = &lunopb.Type{Kind: Object}
var BasicTuple = &lunopb.Type{Kind: Tuple}
var BasicTime = &lunopb.Type{Kind: Time}
var BasicDate = &lunopb.Type{Kind: Date}
var BasicBytes = &lunopb.Type{Kind: Bytes}
var BasicUUID = &lunopb.Type{Kind: UUID}
var BasicTimestamp = &lunopb.Type{Kind: Timestamp}
var BasicInet = &lunopb.Type{Kind: Inet}
var BasicDuration = &lunopb.Type{Kind: Duration}
var BasicRecord = &lunopb.Type{Kind: Record}

// NewTuple constructs a new tuple type containing the given typed items.
func NewTuple(items ...*lunopb.Type) *lunopb.Type {
	tuple := make([]*lunopb.Type, len(items))
	for i, t := range items {
		tuple[i] = (*lunopb.Type)(t)
	}

	return &lunopb.Type{
		Kind:  Tuple,
		Items: tuple,
	}
}

// NewRecord constructs a new record type containing the given typed items.
func NewRecord(items ...*lunopb.Type) *lunopb.Type {
	record := make([]*lunopb.Type, len(items))
	for i, t := range items {
		record[i] = (*lunopb.Type)(t)
	}

	return &lunopb.Type{
		Kind:  Record,
		Items: record,
	}
}

// NewArray constructs a new array type containing the given item type.
func NewArray(item *lunopb.Type) *lunopb.Type {
	return &lunopb.Type{
		Kind:       Array,
		Underlying: (*lunopb.Type)(item),
	}
}
