package scriptbuilder

import (
	"math/big"
)

func (s ScriptBuilder) AllowGas(from, to string, gasPrice, gasLimit *big.Int) ScriptBuilder {
	return s.CallContract("gas", "AllowGas", from, to, gasPrice, gasLimit)
}

func (s ScriptBuilder) SpendGas(address string) ScriptBuilder {
	return s.CallContract("gas", "SpendGas", address)
}

func (s ScriptBuilder) MintTokens(symbol, from, to string, amount *big.Int) ScriptBuilder {
	return s.CallInterop("Runtime.MintTokens", from, to, symbol, amount)
}

func (s ScriptBuilder) Stake(address string, amount *big.Int) ScriptBuilder {
	return s.CallContract("stake", "Stake", address, amount)
}

func (s ScriptBuilder) Unstake(address string, amount *big.Int) ScriptBuilder {
	return s.CallContract("stake", "Unstake", address, amount)
}

func (s ScriptBuilder) TransferTokens(symbol, from, to string, amount *big.Int) ScriptBuilder {
	return s.CallInterop("Runtime.TransferTokens", from, to, symbol, amount)
}

func (s ScriptBuilder) TransferBalance(symbol, from, to string) ScriptBuilder {
	return s.CallInterop("Runtime.TransferTokens", from, to, symbol)
}
