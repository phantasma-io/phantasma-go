package stake

import (
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/domain/types"
	"github.com/phantasma-io/phantasma-go/pkg/io"
)

type EnergyClaim struct {
	StakeAmount *big.Int
	ClaimDate   *types.Timestamp
	IsNew       bool
}

// Serialize implements ther Serializable interface
func (ec *EnergyClaim) Serialize(writer *io.BinWriter) {
	writer.WriteBigInteger(ec.StakeAmount)
	writer.WriteTimestamp(ec.ClaimDate)
	writer.WriteBool(ec.IsNew)
}

// Deserialize implements ther Serializable interface
func (ec *EnergyClaim) Deserialize(reader *io.BinReader) {
	ec.StakeAmount = reader.ReadBigInteger()
	ec.ClaimDate = reader.ReadTimestamp()
	ec.IsNew = reader.ReadBool()
}

type EnergyClaim_S struct {
	StakeAmount string
	ClaimDate   *types.Timestamp
	IsNew       bool
}

// Serialize implements ther Serializable interface
func (ec *EnergyClaim_S) Serialize(writer *io.BinWriter) {
	writer.WriteBigIntegerFromString(ec.StakeAmount)
	writer.WriteTimestamp(ec.ClaimDate)
	writer.WriteBool(ec.IsNew)
}

// Deserialize implements ther Serializable interface
func (ec *EnergyClaim_S) Deserialize(reader *io.BinReader) {
	ec.StakeAmount = reader.ReadBigIntegerToString()
	ec.ClaimDate = reader.ReadTimestamp()
	ec.IsNew = reader.ReadBool()
}
