package contract

import (
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/domain/types"
	"github.com/phantasma-io/phantasma-go/pkg/io"
)

type TokenContent struct {
	SeriesID     *big.Int
	MintID       *big.Int
	CurrentChain string
	Creator      string
	CurrentOwner string
	ROM          []byte
	RAM          []byte
	Timestamp    *types.Timestamp
	Infusion     []TokenInfusion

	// Extra fields, not serializable
	TokenID *big.Int
	Symbol  string
}

func (c *TokenContent) Serialize(writer *io.BinWriter) {
	writer.WriteBigInteger(c.SeriesID)
	writer.WriteBigInteger(c.MintID)

	creator, err := cryptography.FromString(c.Creator)
	if err != nil {
		panic(err)
	}
	creator.Serialize(writer)

	writer.WriteString(c.CurrentChain)

	currentOwner, err := cryptography.FromString(c.CurrentOwner)
	if err != nil {
		panic(err)
	}
	currentOwner.Serialize(writer)

	writer.WriteVarBytes(c.ROM)
	writer.WriteVarBytes(c.RAM)
	writer.WriteTimestamp(c.Timestamp)
	writer.WriteVarUint(uint64(len(c.Infusion)))
	for _, entry := range c.Infusion {
		writer.WriteString(entry.Symbol)
		writer.WriteBigInteger(entry.Value)
	}
}

func (c *TokenContent) Deserialize(reader *io.BinReader) {
	c.SeriesID = reader.ReadBigInteger()
	c.MintID = reader.ReadBigInteger()

	var creator cryptography.Address
	creator.Deserialize(reader)
	c.Creator = creator.String()

	c.CurrentChain = reader.ReadString()

	var currentOwner cryptography.Address
	currentOwner.Deserialize(reader)
	c.CurrentOwner = currentOwner.String()

	c.ROM = reader.ReadVarBytes()

	c.RAM = reader.ReadVarBytes()

	c.Timestamp = reader.ReadTimestamp()

	var infusionCount = reader.ReadVarUint()
	c.Infusion = make([]TokenInfusion, infusionCount)
	for i := 0; i < int(infusionCount); i++ {
		var symbol = reader.ReadString()
		var value = reader.ReadBigInteger()
		c.Infusion[i] = TokenInfusion{Symbol: symbol, Value: value}
	}
}
