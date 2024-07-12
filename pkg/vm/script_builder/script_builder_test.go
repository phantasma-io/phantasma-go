package scriptbuilder_test

import (
	"testing"

	scriptbuilder "github.com/phantasma-io/phantasma-go/pkg/vm/script_builder"
)

func TestNewScript(t *testing.T) {
	sb := scriptbuilder.BeginScript()

	//script := sb.AllowGas("", "", 100000, 800).TransferTokens("SOUL", "SENDER_ADDRESS", "RECEIVER_ADDRESS", 100000).SpendGas().EndScript()
}
