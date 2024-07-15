package scriptbuilder

import (
	"math/big"

	crypto "github.com/phantasma-io/phantasma-go/pkg/cryptography"
)

func (s ScriptBuilder) AllowGas(from, to crypto.Address, gasPrice, gasLimit big.Int) ScriptBuilder {
	return s.CallContract("gas", "AllowGas", from, to, gasPrice, gasLimit)
}

func (s ScriptBuilder) SpendGas(address crypto.Address) ScriptBuilder {
	return s.CallContract("gas", "SpendGas", address)
}

func (s ScriptBuilder) MintTokens(symbol, from, to string, amount big.Int) ScriptBuilder {
	return s.CallInterop("Runtime.MintTokens", from, to, symbol, amount)
}

func (s ScriptBuilder) TransferTokens(symbol, from, to string, amount big.Int) ScriptBuilder {
	return s.CallInterop("Runtime.TransferTokens", from, to, symbol, amount)
}

func (s ScriptBuilder) TransferBalance(symbol, from, to string) ScriptBuilder {
	return s.CallInterop("Runtime.TransferTokens", from, to, symbol)
}
