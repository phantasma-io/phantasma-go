package scriptbuilder

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/phantasma-io/phantasma-go/pkg/io"
	"github.com/phantasma-io/phantasma-go/pkg/vm"
)

type ScriptBuilder struct {
	writer *io.BufBinWriter
}

var _jumpLocations = make(map[int]string)
var _labelLocations = make(map[string]int)

func BeginScript() ScriptBuilder {
	sb := ScriptBuilder{writer: io.NewBufBinWriter()}
	return sb
}

func (s ScriptBuilder) EndScript() []byte {
	s.writer.WriteB(byte(vm.RET))

	return s.writer.Bytes()
}

func (s ScriptBuilder) EmitS(opcode vm.Opcode) ScriptBuilder {
	s.writer.WriteOp(byte(opcode))
	return s
}

func (s ScriptBuilder) EmitM(opcode vm.Opcode, bytes []byte) ScriptBuilder {
	s.writer.WriteOp(byte(opcode))

	if len(bytes) > 0 {
		s.writer.WriteBytes(bytes)
	}

	return s
}

func (s ScriptBuilder) EmitThrow(reg byte) ScriptBuilder {
	s.EmitS(vm.THROW)
	s.writer.WriteB(reg)
	return s
}

func (s ScriptBuilder) EmitPush(reg byte) ScriptBuilder {
	s.EmitS(vm.PUSH)
	s.writer.WriteB(reg)
	return s
}

func (s ScriptBuilder) EmitPop(reg byte) ScriptBuilder {
	s.EmitS(vm.POP)
	s.writer.WriteB(reg)
	return s
}

func (s ScriptBuilder) EmitExtCall(method string, reg byte) ScriptBuilder {
	bytes := []byte(method)
	s.EmitLoad(reg, bytes, vm.String)
	s.EmitS(vm.POP)
	s.writer.WriteB(reg)
	return s
}

func (s ScriptBuilder) EmitLoadBool(reg byte, toLoad bool) ScriptBuilder {

	var bytes []byte
	if toLoad {
		bytes = []byte{1}
	} else {
		bytes = []byte{0}
	}

	s.EmitLoad(reg, bytes, vm.Bool)
	return s
}

func (s ScriptBuilder) EmitLoadTime(reg byte, toLoad time.Time) ScriptBuilder {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(toLoad.Unix()))
	s.EmitLoad(reg, bytes, vm.Timestamp)
	return s
}

func (s ScriptBuilder) EmitLoadInt(reg byte, toLoad int) ScriptBuilder {
	str := strconv.Itoa(toLoad)
	s.EmitLoadString(reg, str)
	return s
}

func (s ScriptBuilder) EmitLoadString(reg byte, toLoad string) ScriptBuilder {
	bytes := []byte(toLoad)
	s.EmitLoad(reg, bytes, vm.String)
	return s
}

func (s ScriptBuilder) EmitLoad(reg byte, bytes []byte, _type vm.VMType) ScriptBuilder {
	// TODO check bytes length
	s.EmitS(vm.LOAD)
	s.writer.WriteB(reg)
	s.writer.WriteB(byte(_type))

	s.writer.WriteVarUint(uint64(len(bytes)))
	s.writer.WriteBytes(bytes)
	return s
}

func (s ScriptBuilder) EmitMove(srcReg byte, dstReg byte) ScriptBuilder {
	s.EmitS(vm.MOVE)
	s.writer.WriteB(srcReg)
	s.writer.WriteB(dstReg)
	return s
}

func (s ScriptBuilder) EmitCopy(srcReg byte, dstReg byte) ScriptBuilder {
	s.EmitS(vm.COPY)
	s.writer.WriteB(srcReg)
	s.writer.WriteB(dstReg)
	return s
}

func (s ScriptBuilder) EmitLabel(label string) ScriptBuilder {
	s.EmitS(vm.NOP)
	// TODO not sure if that works
	_labelLocations[label] = s.writer.Len()
	return s
}

func (s ScriptBuilder) EmitJump(opcode vm.Opcode, label string, reg byte) ScriptBuilder {

	switch opcode {
	case vm.JMP:
	case vm.JMPIF:
	case vm.JMPNOT:
		s.EmitS(opcode)
	default:
		// TODO error
	}

	if opcode != vm.JMP {
		s.writer.WriteB(reg)
	}

	ofs := s.writer.Len()
	s.writer.WriteU16LE(0)
	_jumpLocations[ofs] = label

	return s
}

func (s ScriptBuilder) EmitCall(label string, regCnt byte) ScriptBuilder {
	//TODO register check

	ofs := s.writer.Len() + 2
	s.EmitS(vm.CALL)
	s.writer.WriteB(regCnt)
	s.writer.WriteU16LE(0)

	_jumpLocations[ofs] = label

	return s
}

func (s ScriptBuilder) EmitConditionalJump(opcode vm.Opcode, srcReg byte, label string) ScriptBuilder {
	panic("TODO!!!!!!!!!!!!!!!!")
	//return s
}

func (s ScriptBuilder) EmitVarBytes(value int) ScriptBuilder {
	s.writer.WriteVarUint(uint64(value))
	return s
}

func (s ScriptBuilder) EmitRaw(value []byte) ScriptBuilder {
	s.writer.WriteBytes(value)
	return s
}

func (s ScriptBuilder) loadIntoReg(dstReg byte, arg interface{}) {
	switch e := arg.(type) {
	case string:
		s.EmitLoadString(dstReg, arg.(string))
	case bool:
		s.EmitLoadBool(dstReg, arg.(bool))
	case []byte:
		s.EmitLoad(dstReg, arg.([]byte), vm.Bytes)
	case int:
		s.EmitLoadInt(dstReg, arg.(int))
	case time.Time:
		s.EmitLoadTime(dstReg, arg.(time.Time))
	//TODO array
	default:
		if arg != nil {
			s.writer.Err = errors.New(fmt.Sprintf("unsupported type %s", e))
			return
		}
	}
}

func (s ScriptBuilder) insertMethodArgs(args []interface{}) {
	var tempReg byte = 0

	for i := len(args) - 1; i >= 0; i-- {
		arg := args[i]
		s.loadIntoReg(tempReg, arg)
		s.EmitPush(tempReg)
	}
}

func (s ScriptBuilder) CallInterop(method string, args ...interface{}) ScriptBuilder {
	s.insertMethodArgs(args)

	var dstReg byte = 0
	s.EmitLoadString(dstReg, method)

	s.EmitM(vm.EXTCALL, []byte{dstReg})

	return s
}

func (s ScriptBuilder) CallContract(contractName, method string, args ...interface{}) ScriptBuilder {
	s.insertMethodArgs(args)

	var tmpReg byte = 0
	s.EmitLoadString(tmpReg, method)
	s.EmitPush(tmpReg)

	var srcReg byte = 0
	var dstReg byte = 1
	s.EmitLoadString(srcReg, contractName)
	s.EmitM(vm.CTX, []byte{srcReg, dstReg})

	s.EmitM(vm.SWITCH, []byte{dstReg})
	return s
}
