package contract

import (
	"github.com/phantasma-io/phantasma-go/pkg/io"
	"github.com/phantasma-io/phantasma-go/pkg/vm"
)

type ContractEvent struct {
	Value       byte
	Name        string
	ReturnType  vm.VMType
	Description []byte
}

func (e *ContractEvent) Serialize(writer *io.BinWriter) {
	writer.WriteB(e.Value)
	writer.WriteString(e.Name)
	writer.WriteB(byte(e.ReturnType))
	writer.WriteVarBytes(e.Description)
}

func (e *ContractEvent) Deserialize(reader *io.BinReader) {
	e.Value = reader.ReadB()
	e.Name = reader.ReadString()
	e.ReturnType = vm.VMType(reader.ReadB())
	e.Description = reader.ReadVarBytes()
}
