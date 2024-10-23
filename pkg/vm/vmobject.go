package vm

import (
	"math/big"
	"strconv"

	"github.com/phantasma-io/phantasma-go/pkg/domain/types"
	"github.com/phantasma-io/phantasma-go/pkg/io"
)

type VMObject struct {
	Type VMType
	Data interface{}
}

// AsNumber() returns value stored in vm.VMObject structure, in .Data field, as a *big.Int number
func (v *VMObject) AsNumber() *big.Int {
	switch v.Type {
	case None:
		return big.NewInt(0)

	case String:
		b := big.NewInt(0)
		b.SetString(v.Data.(string), 10)
		return b

	case Bytes:
		// b := v.Data.([]byte)
		// big.NewInt(0).SetBytes(b)
		// TODO probably we need to invert byte order here, not sure if anyone is using this,
		// leaving as unsupported for now
		panic("Not supported!")

	case Enum:
		return big.NewInt(int64(v.Data.(uint32)))

	case Bool:
		var val = v.Data.(bool)
		if val {
			return big.NewInt(1)
		} else {
			return big.NewInt(0)
		}

	case Number:
		n := v.Data.(big.Int)
		return &n

	case Timestamp:
		return big.NewInt(int64(v.Data.(types.Timestamp).Value))

	default:
		panic("Unsupported type")
	}
}

// AsString() returns value stored in vm.VMObject structure, in .Data field, as a string
func (v *VMObject) AsString() string {
	switch v.Type {

	case String:
		return v.Data.(string)

	case Bytes:
		return string(v.Data.([]byte))

	case Enum:
		return strconv.FormatInt(int64(v.Data.(uint32)), 10)

	case Bool:
		var val = v.Data.(bool)
		if val {
			return "true"
		} else {
			return "false"
		}

	case Number:
		n := v.Data.(big.Int)
		return n.String()

	case Timestamp:
		return strconv.FormatInt(int64(v.Data.(types.Timestamp).Value), 10)

	default:
		panic("Unsupported type")
	}
}

// Serialize implements ther Serializable interface
func (v *VMObject) Serialize(writer *io.BinWriter) {
	if v.Type == None {
		return
	}

	writer.WriteB(byte(v.Type))
	switch v.Type {
	case String:
		writer.WriteString(v.Data.(string))
	case Bytes:
		writer.WriteVarBytes(v.Data.([]byte))
	}
}

// Deserialize implements ther Serializable interface
func (v *VMObject) Deserialize(reader *io.BinReader) {
	v.Type = VMType(reader.ReadB())

	switch v.Type {
	case Bool:
		v.Data = reader.ReadBool()
	case Bytes:
		v.Data = reader.ReadVarBytes()
	case Enum:
		v.Data = reader.ReadU32LE()
	case Number:
		v.Data = *reader.ReadBigInteger()
	case Object:
		panic("Not implemented")
	case String:
		v.Data = reader.ReadString()
	case Struct:
		panic("Not implemented")
	case Timestamp:
		v.Data = *reader.ReadTimestamp()
	}
}
