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
