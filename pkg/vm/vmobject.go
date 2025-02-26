package vm

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/phantasma-io/phantasma-go/pkg/domain/types"
	"github.com/phantasma-io/phantasma-go/pkg/io"
	"github.com/phantasma-io/phantasma-go/pkg/util"
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

func (v *VMObject) String() string {
	switch v.Type {
	case None:
		return "Null"
	case Struct:
		return "[Struct]"
	case Bytes:
		return "[Bytes] => " + hex.EncodeToString((v.Data.([]byte)))
	case Number:
		return "[Number] => " + v.AsString()
	case Timestamp:
		return "[Time] => " + v.AsString()
	case String:
		return "[String] => " + v.AsString()
	case Bool:
		return "[Bool] => " + v.AsString()
	case Enum:
		return "[Enum] => " + v.AsString()
	case Object:
		var r string
		if v.Data == nil {
			r = "null"
		} else {
			r = "object"
		}
		return r
	default:
		return "Unknown"
	}
}

func (v *VMObject) SetValue(val []byte, vmtype VMType) *VMObject {
	v.Type = vmtype
	// if val != nil {
	// 	v._localSize = len(val)
	// }

	switch vmtype {
	case Bytes:
		{
			v.Data = val
			break
		}

	case Number:
		{
			var n *big.Int
			if val == nil || len(val) == 0 {
				n = big.NewInt(0)
			} else {
				n = util.BigIntFromCsharpOrPhantasmaByteArray(val)
			}

			v.Data = *n
			break
		}

	case String:
		{
			v.Data = string(val)
			break
		}

	case Enum:
		{
			v.Data = binary.BigEndian.Uint32(val)
			break
		}

	case Timestamp:
		{
			var n uint32
			if val == nil {
				n = 0
			} else {
				n = binary.BigEndian.Uint32(val)
			}
			v.Data = types.Timestamp{Value: n}
			break
		}

	case Bool:
		{
			v.Data = val[0] != 0
			break
		}

	default:
		panic("Unsupported value type")
	}

	return v
}

func (v *VMObject) Copy(other *VMObject) {
	if other == nil || other.Type == None {
		v.Type = None
		v.Data = nil
		return
	}

	v.Type = other.Type

	if other.Type == Struct {
		panic("Struct type copying is unsupported")
	} else {
		v.Data = other.Data
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

func ValidateStructKey(key *VMObject) {
	if key.Type == None {
		panic("Cannot use value of type None as key for struct field")
	}
	if key.Type == Struct {
		panic("Cannot use value of type Struct as key for struct field")
	}
	if key.Type == Object {
		panic("Cannot use value of type Object as key for struct field")
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
		childCount := reader.ReadU32LE()
		children := make(map[VMObject]VMObject)
		for {
			if childCount == 0 {
				break
			}

			key := &VMObject{}
			key.Deserialize(reader)

			ValidateStructKey(key)

			val := &VMObject{}
			val.Deserialize(reader)

			children[*key] = *val
			childCount--
		}

		v.Data = children
	case Timestamp:
		v.Data = *reader.ReadTimestamp()
	}
}
