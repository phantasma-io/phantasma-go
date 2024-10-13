package rpc

import (
	"fmt"
	"math/big"

	"context"

	"github.com/phantasma-io/phantasma-go/pkg/jsonrpc"
	resp "github.com/phantasma-io/phantasma-go/pkg/rpc/response"
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

// NewRPCSetMainnet returns a new set of RPC clients
func NewRPCSetMainnet() []PhantasmaRPC {
	return []PhantasmaRPC{NewRPC("https://pharpc1.phantasma.info/rpc"),
		NewRPC("https://pharpc2.phantasma.info/rpc"),
		NewRPC("https://pharpc3.phantasma.info/rpc"),
		NewRPC("https://pharpc4.phantasma.info/rpc")}
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
	result, err := rpc.client.Call(context.Background(), "getPlatforms", nil)

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
	result, err := rpc.client.Call(context.Background(), "getAccounts", addresses, false)

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
	result, err := rpc.client.Call(context.Background(), "getAccount", address, false)

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
	result, err := rpc.client.Call(context.Background(), "getAccount", address, false)

	if err := checkError(err, result.Error); err != nil {
		return resp.AccountResult{}, err
	}

	err = result.GetObject(&account)
	if err != nil {
		return resp.AccountResult{}, err
	}

	return account, nil
}

// Deprecated: Long execution time and possibility of inconsistent result
// GetAccountEx returns current account state and list of all txes, including latest tx for this account (last in tx list)
func (rpc PhantasmaRPC) GetAccountEx(address string) (resp.AccountResult, error) {
	var account resp.AccountResult
	result, err := rpc.client.Call(context.Background(), "getAccount", address, true)

	if err := checkError(err, result.Error); err != nil {
		return resp.AccountResult{}, err
	}

	err = result.GetObject(&account)
	if err != nil {
		return resp.AccountResult{}, err
	}

	return account, nil
}

// GetAddressTransactions Returns list of transactions for given address
// Transactions are ordered from newer to older
func (rpc PhantasmaRPC) GetAddressTransactions(address string, page int, pageSize int) (resp.PaginatedResult[resp.AddressTransactionsResult], error) {
	var addressTxs resp.PaginatedResult[resp.AddressTransactionsResult]
	result, err := rpc.client.Call(context.Background(), "getAddressTransactions", address, page, pageSize)

	if err := checkError(err, result.Error); err != nil {
		return resp.PaginatedResult[resp.AddressTransactionsResult]{}, err
	}

	err = result.GetObject(&addressTxs)
	if err != nil {
		return resp.PaginatedResult[resp.AddressTransactionsResult]{}, err
	}

	return addressTxs, nil
}

// GetAddressTransactionCount Returns number of transactions for given address
func (rpc PhantasmaRPC) GetAddressTransactionCount(address string, chainName string) (int, error) {
	var count int
	result, err := rpc.client.Call(context.Background(), "getAddressTransactionCount", address, chainName)

	if err := checkError(err, result.Error); err != nil {
		return 0, err
	}

	err = result.GetObject(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetBlockByHeight Returns block by height
func (rpc PhantasmaRPC) GetBlockByHeight(chain string, height string) (resp.BlockResult, error) {
	var blockResult resp.BlockResult
	result, err := rpc.client.Call(context.Background(), "getBlockByHeight", chain, height)

	if err := checkError(err, result.Error); err != nil {
		return resp.BlockResult{}, err
	}

	err = result.GetObject(&blockResult)
	if err != nil {
		errorResult := resp.ErrorResult{}
		err = result.GetObject(&errorResult)
		if err != nil {
			return blockResult, err
		}

		return blockResult, fmt.Errorf(errorResult.Error)
	}
	return blockResult, nil
}

// GetBlockHeight Returns height of the latest block minted on the chain
func (rpc PhantasmaRPC) GetBlockHeight(chainName string) (*big.Int, error) {
	var resultValue string
	result, err := rpc.client.Call(context.Background(), "getBlockHeight", chainName)

	if err := checkError(err, result.Error); err != nil {
		return big.NewInt(0), err
	}

	err = result.GetObject(&resultValue)
	if err != nil {
		return big.NewInt(0), err
	}

	height, _ := big.NewInt(0).SetString(resultValue, 10)
	return height, nil
}

// InvokeRawScript comment
func (rpc PhantasmaRPC) InvokeRawScript(chain, script string) (resp.ScriptResult, error) {
	scriptResult := resp.ScriptResult{}
	result, err := rpc.client.Call(context.Background(), "invokeRawScript", chain, script)

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
	result, err := rpc.client.Call(context.Background(), "sendRawTransaction", txData)

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
	result, err := rpc.client.Call(context.Background(), "getTransaction", txHash)

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
	result, err := rpc.client.Call(context.Background(), "getTokens", extended)

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
	result, err := rpc.client.Call(context.Background(), "getToken", symbol, extended)

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
