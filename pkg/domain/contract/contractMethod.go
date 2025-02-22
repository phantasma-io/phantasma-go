package contract

import (
	"github.com/phantasma-io/phantasma-go/pkg/io"
	"github.com/phantasma-io/phantasma-go/pkg/vm"
)

type ContractMethod struct {
	Name       string
	ReturnType vm.VMType
	Parameters []ContractParameter
	Offset     int32
}

func (m *ContractMethod) Serialize(writer *io.BinWriter) {
	writer.WriteString(m.Name)
	writer.WriteB(byte(m.ReturnType))
	writer.WriteU32LE(uint32(m.Offset))
	writer.WriteB(byte(len(m.Parameters)))
	for _, entry := range m.Parameters {
		writer.WriteString(entry.Name)
		writer.WriteB(byte(entry.Type))
	}
}

func (m *ContractMethod) Deserialize(reader *io.BinReader) {
	m.Name = reader.ReadString()
	m.ReturnType = vm.VMType(reader.ReadB())
	m.Offset = int32(reader.ReadU32LE())
	len := int(reader.ReadB())
	m.Parameters = make([]ContractParameter, len)
	for i := 0; i < len; i++ {
		var pName = reader.ReadString()
		var pVMType = vm.VMType(reader.ReadB())
		m.Parameters[i] = ContractParameter{Name: pName, Type: pVMType}
	}
}

type ContractMethod_S struct {
	Name       string
	ReturnType string
	Parameters []ContractParameter
	Offset     int32
}

func (m *ContractMethod_S) Serialize(writer *io.BinWriter) {
	writer.WriteString(m.Name)

	var vmType vm.VMType
	writer.WriteB(byte(vmType.FromString(m.ReturnType)))

	writer.WriteU32LE(uint32(m.Offset))
	writer.WriteB(byte(len(m.Parameters)))
	for _, entry := range m.Parameters {
		writer.WriteString(entry.Name)
		writer.WriteB(byte(entry.Type))
	}
}

func (m *ContractMethod_S) Deserialize(reader *io.BinReader) {
	m.Name = reader.ReadString()

	m.ReturnType = vm.VMTypeLookup[vm.VMType(reader.ReadB())]

	m.Offset = int32(reader.ReadU32LE())
	len := int(reader.ReadB())
	m.Parameters = make([]ContractParameter, len)
	for i := 0; i < len; i++ {
		var pName = reader.ReadString()
		var pVMType = vm.VMType(reader.ReadB())
		m.Parameters[i] = ContractParameter{Name: pName, Type: pVMType}
	}
}
