package scriptbuilder_test

import (
	"math/big"
	"testing"

	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
	scriptbuilder "github.com/phantasma-io/phantasma-go/pkg/vm/script_builder"
	"github.com/stretchr/testify/assert"
)

func TestNewScript(t *testing.T) {
	fromAddress, _ := cryptography.FromString("P2KM9FjYrDXnPPAynLXAHdQ8wYz8de9VbDeybrLepnw6C5x")
	toAddress, _ := cryptography.FromString("P2KM9FjYrDXnPPAynLXAHdQ8wYz8de9VbDeybrLepnw6C5x")
	symbols := []string{"SOUL", "KCAL"}

	assert.NotPanics(t, func() {
		sb := scriptbuilder.BeginScript()
		sb.AllowGas(fromAddress, cryptography.NullAddress(), big.NewInt(100000), big.NewInt(21000)).
			TransferTokens(symbols[0], fromAddress, toAddress, big.NewInt(100000000)).
			SpendGas(fromAddress).
			EndScript()
	})

	assert.NotPanics(t, func() {
		sb := scriptbuilder.BeginScript()
		sb.CallInterop("Runtime.TransferToken", fromAddress, toAddress, symbols[0], "TOKEN_ID")
	})

	assert.Panics(t, func() {
		sb := scriptbuilder.BeginScript()
		// Arrays are not supported by script builder
		sb.CallInterop("Runtime.TransferToken", fromAddress, toAddress, symbols, "TOKEN_ID")
	})
}
