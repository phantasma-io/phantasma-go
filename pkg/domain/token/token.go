package token

import (
	"math/big"
	"slices"

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

var tokenFlagsLookup = map[TokenFlags]string{
	None:         `None`,
	Transferable: `Transferable`,
	Fungible:     `Fungible`,
	Finite:       `Finite`,
	Divisible:    `Divisible`,
	Fuel:         `Fuel`,
	Stakable:     `Stakable`,
	Fiat:         `Fiat`,
	Swappable:    `Swappable`,
	Burnable:     `Burnable`,
	Mintable:     `Mintable`,
}

func (tf TokenFlags) hasFlag(flagToCheck TokenFlags) bool {
	val := tf & flagToCheck
	return (val > 0)
}

func (tf TokenFlags) setFlag(flagToSet TokenFlags) TokenFlags {
	tf |= flagToSet
	return tf
}

func (tf TokenFlags) appendStringFlagIfSet(flags []string, flagToCheck TokenFlags) []string {
	if tf.hasFlag(flagToCheck) {
		flags = append(flags, tokenFlagsLookup[flagToCheck])
	}

	return flags
}

func (tf TokenFlags) appendFlagIfSet(flags []string, flagToCheck TokenFlags) TokenFlags {
	if slices.Contains(flags, tokenFlagsLookup[flagToCheck]) {
		tf = tf.setFlag(flagToCheck)
	}

	return tf
}

func (tf TokenFlags) FromSlice(flagsSlice []string) TokenFlags {
	tf = None
	tf = tf.appendFlagIfSet(flagsSlice, Transferable)
	tf = tf.appendFlagIfSet(flagsSlice, Fungible)
	tf = tf.appendFlagIfSet(flagsSlice, Finite)
	tf = tf.appendFlagIfSet(flagsSlice, Divisible)
	tf = tf.appendFlagIfSet(flagsSlice, Fuel)
	tf = tf.appendFlagIfSet(flagsSlice, Stakable)
	tf = tf.appendFlagIfSet(flagsSlice, Fiat)
	tf = tf.appendFlagIfSet(flagsSlice, Swappable)
	tf = tf.appendFlagIfSet(flagsSlice, Burnable)
	tf = tf.appendFlagIfSet(flagsSlice, Mintable)

	return tf
}

func (tf TokenFlags) ToSlice() []string {
	result := make([]string, 0)
	result = tf.appendStringFlagIfSet(result, Transferable)
	result = tf.appendStringFlagIfSet(result, Fungible)
	result = tf.appendStringFlagIfSet(result, Finite)
	result = tf.appendStringFlagIfSet(result, Divisible)
	result = tf.appendStringFlagIfSet(result, Fuel)
	result = tf.appendStringFlagIfSet(result, Stakable)
	result = tf.appendStringFlagIfSet(result, Fiat)
	result = tf.appendStringFlagIfSet(result, Swappable)
	result = tf.appendStringFlagIfSet(result, Burnable)
	result = tf.appendStringFlagIfSet(result, Mintable)

	if len(result) == 0 {
		result = append(result, tokenFlagsLookup[None])
	}

	return result
}

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

type TokenInfo_S struct {
	Symbol    string
	Name      string
	Owner     cryptography.Address
	Flags     []string
	MaxSupply string
	Decimals  int32
	Script    []byte
	ABI       contract.ContractInterface_S
}

func (ti *TokenInfo_S) Serialize(writer *io.BinWriter) {
	writer.WriteString(ti.Symbol)
	writer.WriteString(ti.Name)
	ti.Owner.Serialize(writer)

	var tf TokenFlags
	writer.WriteU32LE(uint32(tf.FromSlice(ti.Flags)))

	writer.WriteU32LE(uint32(ti.Decimals))
	writer.WriteBigIntegerFromString(ti.MaxSupply)
	writer.WriteVarBytes(ti.Script)

	bytes := io.Serialize[*contract.ContractInterface_S](&ti.ABI)
	writer.WriteVarBytes(bytes)
}

func (ti *TokenInfo_S) Deserialize(reader *io.BinReader) {
	ti.Symbol = reader.ReadString()
	ti.Name = reader.ReadString()
	ti.Owner.Deserialize(reader)

	ti.Flags = TokenFlags(reader.ReadU32LE()).ToSlice()

	ti.Decimals = int32(reader.ReadU32LE())
	ti.MaxSupply = reader.ReadBigIntegerToString()
	ti.Script = reader.ReadVarBytes()

	bytes := reader.ReadVarBytes()
	ti.ABI = *io.Deserialize[*contract.ContractInterface_S](bytes)
}
