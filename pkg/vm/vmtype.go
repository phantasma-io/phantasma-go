package vm

//go:generate stringer -type=Opcode -linecomment

// VMType identifies the type of a vm object
type VMType byte

// Viable list of supported instruction constants.
const (
	None VMType = iota
	Struct
	Bytes
	Number
	String
	Timestamp
	Bool
	Enum
	Object
)

var VMTypeLookup = map[VMType]string{
	None:      `None`,
	Struct:    `Struct`,
	Bytes:     `Bytes`,
	Number:    `Number`,
	String:    `String`,
	Timestamp: `Timestamp`,
	Bool:      `Bool`,
	Enum:      `Enum`,
	Object:    `Object`,
}

func (t VMType) FromString(vmType string) VMType {
	for k1, s := range VMTypeLookup {
		if s == vmType {

			return k1
		}
	}

	return None
}
