package scriptbuilder

import (
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
)

func (s ScriptBuilder) AllowGas(from, to cryptography.Address, gasPrice, gasLimit *big.Int) ScriptBuilder {
	return s.CallContract("gas", "AllowGas", from, to, gasPrice, gasLimit)
}

func (s ScriptBuilder) SpendGas(address cryptography.Address) ScriptBuilder {
	return s.CallContract("gas", "SpendGas", address)
}

func (s ScriptBuilder) MintTokens(symbol string, from, to cryptography.Address, amount *big.Int) ScriptBuilder {
	return s.CallInterop("Runtime.MintTokens", from, to, symbol, amount)
}

func (s ScriptBuilder) Stake(address cryptography.Address, amount *big.Int) ScriptBuilder {
	return s.CallContract("stake", "Stake", address, amount)
}

func (s ScriptBuilder) Unstake(address cryptography.Address, amount *big.Int) ScriptBuilder {
	return s.CallContract("stake", "Unstake", address, amount)
}

func (s ScriptBuilder) TransferTokens(symbol string, from, to cryptography.Address, amount *big.Int) ScriptBuilder {
	return s.CallInterop("Runtime.TransferTokens", from, to, symbol, amount)
}

func (s ScriptBuilder) TransferBalance(symbol string, from, to cryptography.Address) ScriptBuilder {
	return s.CallInterop("Runtime.TransferTokens", from, to, symbol)
}
