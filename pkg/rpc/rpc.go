package rpc

import (
	"fmt"

	resp "github.com/phantasma-io/phantasma-go/pkg/rpc/response"
	"github.com/ybbus/jsonrpc/v2"
)

// PhantasmaRPC struct
type PhantasmaRPC struct {
	client jsonrpc.RPCClient
}

// TODO provide multiple clients with fallback functionality and penalty records on each client
// if a client fails to respond, increas penalty counter, making the possibility to be choosen less
// likely.

// NewRPCs returns a new RPC client
//func NewRPCs(endpoints []string) PhantasmaRPC {
//	rpc := PhantasmaRPC{
//		client: jsonrpc.NewClient("http://207.148.17.86:7077/rpc"),
//	}
//}

// NewRPCMainnet returns a new RPC client
func NewRPCMainnet() PhantasmaRPC {
	return NewRPC("https://pharpc1.phantasma.info/rpc")
}

// NewRPCTestnet returns a new testnet RPC client
func NewRPCTestnet() PhantasmaRPC {
	return NewRPC("https://testnet.phantasma.info/rpc")
}

// NewRPC returns a new RPC client
func NewRPC(endpoint string) PhantasmaRPC {
	rpc := PhantasmaRPC{
		client: jsonrpc.NewClient(endpoint),
	}
	return rpc
}

func checkError(err error, rpcError *jsonrpc.RPCError) error {
	if err != nil {
		return err
	}

	if rpcError != nil {
		return fmt.Errorf(rpcError.Message)
	}

	return nil
}

// GetPlatforms comment
func (rpc PhantasmaRPC) GetPlatforms() ([]resp.PlatformResult, error) {
	var platforms []resp.PlatformResult
	result, err := rpc.client.Call("getPlatforms", nil)

	if err := checkError(err, result.Error); err != nil {
		return []resp.PlatformResult{}, err
	}

	err = result.GetObject(&platforms)
	if err != nil {
		return []resp.PlatformResult{}, err
	}

	return platforms, nil
}

// GetAccounts takes a comma separated list of addresses
func (rpc PhantasmaRPC) GetAccounts(addresses string) ([]resp.AccountResult, error) {
	var accounts []resp.AccountResult
	result, err := rpc.client.Call("getAccounts", addresses)

	if err := checkError(err, result.Error); err != nil {
		return []resp.AccountResult{}, err
	}

	err = result.GetObject(&accounts)
	if err != nil {
		return []resp.AccountResult{}, err
	}

	return accounts, nil
}

// LookupName comment
func (rpc PhantasmaRPC) LookupName(name string) (string, error) {
	var address string
	result, err := rpc.client.Call("getAccount", address)

	if err := checkError(err, result.Error); err != nil {
		return "", err
	}

	address, err = result.GetString()
	if err != nil {
		return "", err
	}

	return name, nil
}

// GetAccount comment
func (rpc PhantasmaRPC) GetAccount(address string) (resp.AccountResult, error) {
	var account resp.AccountResult
	result, err := rpc.client.Call("getAccount", address)

	if err := checkError(err, result.Error); err != nil {
		return resp.AccountResult{}, err
	}

	err = result.GetObject(&account)
	if err != nil {
		return resp.AccountResult{}, err
	}

	return account, nil
}

// InvokeRawScript comment
func (rpc PhantasmaRPC) InvokeRawScript(chain, script string) (resp.ScriptResult, error) {
	scriptResult := resp.ScriptResult{}
	result, err := rpc.client.Call("invokeRawScript", chain, script)

	if err := checkError(err, result.Error); err != nil {
		return resp.ScriptResult{}, err
	}

	err = result.GetObject(&scriptResult)
	if err != nil {
		errorResult := resp.ErrorResult{}
		err = result.GetObject(&errorResult)
		if err != nil {
			return scriptResult, err
		}

		return scriptResult, fmt.Errorf(errorResult.Error)
	}

	return scriptResult, nil
}

// SendRawTransaction comment
func (rpc PhantasmaRPC) SendRawTransaction(txData string) (string, error) {
	var hash string
	result, err := rpc.client.Call("sendRawTransaction", txData)

	if err := checkError(err, result.Error); err != nil {
		return "", err
	}

	hash, err = result.GetString()
	if err != nil {
		errorResult := resp.ErrorResult{}
		err = result.GetObject(&errorResult)
		if err != nil {
			return hash, err
		}

		return hash, fmt.Errorf(errorResult.Error)
	}

	return hash, nil
}

// GetTransaction comment
func (rpc PhantasmaRPC) GetTransaction(txHash string) (resp.TransactionResult, error) {
	var txResult resp.TransactionResult
	result, err := rpc.client.Call("getTransaction", txHash)

	if err := checkError(err, result.Error); err != nil {
		return resp.TransactionResult{}, err
	}

	err = result.GetObject(&txResult)
	if err != nil {
		errorResult := resp.ErrorResult{}
		err = result.GetObject(&errorResult)
		if err != nil {
			return txResult, err
		}

		return txResult, fmt.Errorf(errorResult.Error)
	}
	return txResult, nil
}

// GetTokens comment
func (rpc PhantasmaRPC) GetTokens(extended bool) ([]resp.TokenResult, error) {
	var txResult []resp.TokenResult
	result, err := rpc.client.Call("getTokens", extended)

	if err := checkError(err, result.Error); err != nil {
		return []resp.TokenResult{}, err
	}

	err = result.GetObject(&txResult)
	if err != nil {
		errorResult := resp.ErrorResult{}
		err = result.GetObject(&errorResult)
		if err != nil {
			return txResult, err
		}

		return txResult, fmt.Errorf(errorResult.Error)
	}
	return txResult, nil
}

// GetToken comment
func (rpc PhantasmaRPC) GetToken(symbol string, extended bool) (resp.TokenResult, error) {
	var txResult resp.TokenResult
	result, err := rpc.client.Call("getToken", symbol, extended)

	if err := checkError(err, result.Error); err != nil {
		return resp.TokenResult{}, err
	}

	err = result.GetObject(&txResult)
	if err != nil {
		errorResult := resp.ErrorResult{}
		err = result.GetObject(&errorResult)
		if err != nil {
			return txResult, err
		}

		return txResult, fmt.Errorf(errorResult.Error)
	}
	return txResult, nil
}
