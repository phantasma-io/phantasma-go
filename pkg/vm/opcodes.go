package vm

//go:generate stringer -type=Opcode -linecomment

// Opcode represents a single operation code for the Phantasma virtual machine.
type Opcode byte

// Viable list of supported instruction constants.
const (
	NOP Opcode = iota

	// register
	MOVE // copy reference
	COPY // copy by value
	PUSH
	POP
	SWAP

	// flow
	CALL
	EXTCALL
	JMP
	JMPIF
	JMPNOT
	RET
	THROW

	// data
	LOAD
	CAST
	CAT
	RANGE
	LEFT
	RIGHT
	SIZE
	COUNT

	// logical
	NOT
	AND
	OR
	XOR
	EQUAL
	LT
	GT
	LTE
	GTE

	// numeric
	INC
	DEC
	SIGN
	NEGATE
	ABS
	ADD
	SUB
	MUL
	DIV
	MOD
	SHL
	SHR
	MIN
	MAX
	POW

	// context
	CTX
	SWITCH

	// array
	PUT
	GET    // lookups a key and copies a reference into register
	CLEAR  // clears a register
	UNPACK // unpacks serialized struct based on ref struct
	PACK   // unused for now

	//  debugger
	DEBUG
)
