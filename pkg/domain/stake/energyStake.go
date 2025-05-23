package stake

import (
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/domain/types"
	"github.com/phantasma-io/phantasma-go/pkg/io"
)

type EnergyStake struct {
	StakeAmount *big.Int
	StakeTime   *types.Timestamp
}

// Serialize implements ther Serializable interface
func (es *EnergyStake) Serialize(writer *io.BinWriter) {
	writer.WriteBigInteger(es.StakeAmount)
	writer.WriteTimestamp(es.StakeTime)
}

// Deserialize implements ther Serializable interface
func (es *EnergyStake) Deserialize(reader *io.BinReader) {
	es.StakeAmount = reader.ReadBigInteger()
	es.StakeTime = reader.ReadTimestamp()
}

type EnergyStake_S struct {
	StakeAmount string
	StakeTime   *types.Timestamp
}

// Serialize implements ther Serializable interface
func (es *EnergyStake_S) Serialize(writer *io.BinWriter) {
	writer.WriteBigIntegerFromString(es.StakeAmount)
	writer.WriteTimestamp(es.StakeTime)
}

// Deserialize implements ther Serializable interface
func (es *EnergyStake_S) Deserialize(reader *io.BinReader) {
	es.StakeAmount = reader.ReadBigIntegerToString()
	es.StakeTime = reader.ReadTimestamp()
}
