package contract

import (
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/io"
)

type TokenSeriesMode uint

const (
	Unique     TokenSeriesMode = 1
	Duplicated TokenSeriesMode = 2
)

type TokenSeries struct {
	MintCount *big.Int
	MaxSupply *big.Int
	Mode      TokenSeriesMode
	Script    []byte
	ABI       ContractInterface
	ROM       []byte
}

func (s *TokenSeries) Serialize(writer *io.BinWriter) {

	writer.WriteBigInteger(s.MintCount)
	writer.WriteBigInteger(s.MaxSupply)
	writer.WriteB(byte(s.Mode))
	writer.WriteVarBytes(s.Script)

	// TODO this is how it should have been done for ABI
	// But it's incompatible with current storage
	// s.ABI.Serialize(writer)

	bytes := io.Serialize[*ContractInterface](&s.ABI)
	writer.WriteVarBytes(bytes)

	writer.WriteVarBytes(s.ROM)
}

func (s *TokenSeries) Deserialize(reader *io.BinReader) {
	s.MintCount = reader.ReadBigInteger()
	s.MaxSupply = reader.ReadBigInteger()
	s.Mode = TokenSeriesMode(reader.ReadB())
	s.Script = reader.ReadVarBytes()

	// TODO this is how it should have been done for ABI
	// This length byte was stored but we don't really need it
	// s.ABI = ContractInterface{}
	// s.ABI.Deserialize(reader)

	bytes := reader.ReadVarBytes()
	s.ABI = *io.Deserialize[*ContractInterface](bytes)

	s.ROM = reader.ReadVarBytes()
}
