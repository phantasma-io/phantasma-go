package io

import (
	"bytes"
)

// Serializer defines the binary encoding/decoding interface. Errors are
// returned via BinReader/BinWriter Err field. These functions must have safe
// behavior when passed BinReader/BinWriter with Err already set. Invocations
// to these functions tend to be nested, with this mechanism only the top-level
// caller should handle the error once and all the other code should just not
// panic in presence of error.
type Serializer interface {
	Deserialize(*BinReader)
	Serialize(*BinWriter)
}

type deserializable interface {
	Deserialize(*BinReader)
}

type serializable interface {
	Serialize(*BinWriter)
}

func Deserialize[T deserializable](data []byte, object T) T {
	br := NewBinReaderFromBuf(data)

	object.Deserialize(br)
	return object
}

func Serialize[T serializable](object T) []byte {
	b := new(bytes.Buffer)
	bw := NewBinWriterFromIO(b)

	object.Serialize(bw)

	if bw.Err != nil {
		return nil
	}

	return b.Bytes()
}

// Generic serializer/deserializer - unfished, not sure if worth finishing
//// SerializeType comment
//func SerializeType(writer *BinWriter, i interface{}) {
//	switch v := i.(type) {
//	case bool:
//		if v {
//			writer.WriteB(byte(1))
//		} else {
//			writer.WriteB(byte(0))
//		}
//	case byte:
//		writer.WriteB(v)
//	case int32:
//		writer.WriteU32LE(uint32(v))
//	case int64:
//		writer.WriteU64LE(uint64(v))
//	case uint32:
//		writer.WriteU32LE(v)
//	case uint64:
//		writer.WriteU64LE(v)
//	case string:
//		writer.WriteString(v)
//	case big.Int:
//		writer.WriteBytes(v.Bytes())
//	case Serializer:
//		v.Serialize(writer)
//	case []interface{}:
//		if v == nil {
//			writer.WriteVarUint(0)
//		}
//
//		for _, elem := range v {
//			SerializeType(writer, elem)
//		}
//	case interface{}:
//		t := reflect.ValueOf(v)
//		for i := 0; i < t.NumField(); i++ {
//			val := t.Field(i)
//			SerializeType(writer, val)
//		}
//	}
//}
//
//// DeserializeType comment
//func DeserializeType(reader *BinReader, i interface{}) interface{} {
//	switch v := i.(type) {
//	case bool:
//		return reader.ReadB() != 0
//	case byte:
//		return reader.ReadB()
//	case int32:
//		reader.ReadU32LE()
//	case int64:
//		reader.ReadU64LE()
//	case uint32:
//		reader.ReadU32LE()
//	case uint64:
//		reader.ReadU64LE()
//	case string:
//		reader.ReadString()
//	case big.Int:
//		reader.ReadBytes(v.Bytes())
//	//case Serializer:
//	//	return v.Deserialize(reader)
//	case []interface{}:
//		if v == nil {
//			reader.ReadVarUint()
//		}
//
//		for _, elem := range v {
//			return DeserializeType(reader, elem)
//		}
//	case interface{}:
//		t := reflect.ValueOf(v)
//		for i := 0; i < t.NumField(); i++ {
//			val := t.Field(i)
//			SerializeType(writer, val)
//		}
//	}
//}
