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
