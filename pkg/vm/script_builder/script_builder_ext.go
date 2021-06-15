package scriptbuilder

func (s ScriptBuilder) AllowGas(from, to string, gasPrice, gasLimit int) ScriptBuilder {
	return s.CallContract("gas", "AllowGas", from, to, gasPrice, gasLimit)
}

func (s ScriptBuilder) SpendGas(address string) ScriptBuilder {
	return s.CallContract("gas", "SpendGas", address)
}

func (s ScriptBuilder) MintTokens(symbol, from, to string, amount int) ScriptBuilder {
	return s.CallInterop("Runtime.MintTokens", from, to, symbol, amount)
}

func (s ScriptBuilder) TransferTokens(symbol, from, to string, amount int) ScriptBuilder {
	return s.CallInterop("Runtime.TransferTokens", from, to, symbol, amount)
}

func (s ScriptBuilder) TransferBalance(symbol, from, to string) ScriptBuilder {
	return s.CallInterop("Runtime.TransferTokens", from, to, symbol)
}
