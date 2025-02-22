package contract

import (
	"github.com/phantasma-io/phantasma-go/pkg/io"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type ContractInterface struct {
	Methods *orderedmap.OrderedMap[string, ContractMethod]

	Events []ContractEvent
}

func (i *ContractInterface) Serialize(writer *io.BinWriter) {
	{
		writer.WriteB(byte(i.Methods.Len()))
		for pair := i.Methods.Oldest(); pair != nil; pair = pair.Next() {
			pair.Value.Serialize(writer)
		}

		writer.WriteB(byte(len(i.Events)))
		for _, evt := range i.Events {
			evt.Serialize(writer)
		}
	}
}

func (i *ContractInterface) Deserialize(reader *io.BinReader) {
	len := int(reader.ReadB())
	i.Methods = orderedmap.New[string, ContractMethod]()
	for j := 0; j < len; j++ {
		method := ContractMethod{}
		method.Deserialize(reader)
		i.Methods.Set(method.Name, method)
	}

	len = int(reader.ReadB())
	i.Events = make([]ContractEvent, len)
	for j := 0; j < len; j++ {
		e := ContractEvent{}
		e.Deserialize(reader)
		i.Events[j] = e
	}
}

type ContractInterface_S struct {
	Methods *orderedmap.OrderedMap[string, ContractMethod_S]

	Events []ContractEvent
}

func (i *ContractInterface_S) Serialize(writer *io.BinWriter) {
	{
		writer.WriteB(byte(i.Methods.Len()))
		for pair := i.Methods.Oldest(); pair != nil; pair = pair.Next() {
			pair.Value.Serialize(writer)
		}

		writer.WriteB(byte(len(i.Events)))
		for _, evt := range i.Events {
			evt.Serialize(writer)
		}
	}
}

func (i *ContractInterface_S) Deserialize(reader *io.BinReader) {
	len := int(reader.ReadB())
	i.Methods = orderedmap.New[string, ContractMethod_S]()
	for j := 0; j < len; j++ {
		method := ContractMethod_S{}
		method.Deserialize(reader)
		i.Methods.Set(method.Name, method)
	}

	len = int(reader.ReadB())
	i.Events = make([]ContractEvent, len)
	for j := 0; j < len; j++ {
		e := ContractEvent{}
		e.Deserialize(reader)
		i.Events[j] = e
	}
}
