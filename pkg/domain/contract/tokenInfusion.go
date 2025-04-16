package contract

import (
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/io"
)

type TokenInfusion struct {
	Symbol string
	Value  *big.Int
}

func (i *TokenInfusion) Serialize(writer *io.BinWriter) {
	writer.WriteString(i.Symbol)
	writer.WriteBigInteger(i.Value)
}

func (i *TokenInfusion) Deserialize(reader *io.BinReader) {
	i.Symbol = reader.ReadString()
	i.Value = reader.ReadBigInteger()
}

type TokenInfusion_S struct {
	Symbol string
	Value  string
}

func (i *TokenInfusion_S) Serialize(writer *io.BinWriter) {
	writer.WriteString(i.Symbol)
	writer.WriteBigIntegerFromString(i.Value)
}

func (i *TokenInfusion_S) Deserialize(reader *io.BinReader) {
	i.Symbol = reader.ReadString()
	i.Value = reader.ReadBigIntegerToString()
}
