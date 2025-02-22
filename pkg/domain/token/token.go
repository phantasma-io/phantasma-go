package token

import (
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/domain/contract"
	"github.com/phantasma-io/phantasma-go/pkg/io"
)

type TokenFlags uint

const (
	None         TokenFlags = 0
	Transferable TokenFlags = 1 << 0
	Fungible     TokenFlags = 1 << 1
	Finite       TokenFlags = 1 << 2
	Divisible    TokenFlags = 1 << 3
	Fuel         TokenFlags = 1 << 4
	Stakable     TokenFlags = 1 << 5
	Fiat         TokenFlags = 1 << 6
	Swappable    TokenFlags = 1 << 7
	Burnable     TokenFlags = 1 << 8
	Mintable     TokenFlags = 1 << 9
)

type TokenInfo struct {
	Symbol    string
	Name      string
	Owner     cryptography.Address
	Flags     TokenFlags
	MaxSupply *big.Int
	Decimals  int32
	Script    []byte
	ABI       contract.ContractInterface
}

func (ti *TokenInfo) Serialize(writer *io.BinWriter) {
	writer.WriteString(ti.Symbol)
	writer.WriteString(ti.Name)
	ti.Owner.Serialize(writer)
	writer.WriteU32LE(uint32(ti.Flags))
	writer.WriteU32LE(uint32(ti.Decimals))
	writer.WriteBigInteger(ti.MaxSupply)
	writer.WriteVarBytes(ti.Script)

	bytes := io.Serialize[*contract.ContractInterface](&ti.ABI)
	writer.WriteVarBytes(bytes)
}

func (ti *TokenInfo) Deserialize(reader *io.BinReader) {
	ti.Symbol = reader.ReadString()
	ti.Name = reader.ReadString()
	ti.Owner.Deserialize(reader)
	ti.Flags = TokenFlags(reader.ReadU32LE())
	ti.Decimals = int32(reader.ReadU32LE())
	ti.MaxSupply = reader.ReadBigInteger()
	ti.Script = reader.ReadVarBytes()

	bytes := reader.ReadVarBytes()
	ti.ABI = *io.Deserialize[*contract.ContractInterface](bytes)
}
