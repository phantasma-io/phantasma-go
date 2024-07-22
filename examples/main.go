package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	chain "github.com/phantasma-io/phantasma-go/pkg/blockchain"
	crypto "github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/rpc"
	"github.com/phantasma-io/phantasma-go/pkg/rpc/response"
	"github.com/phantasma-io/phantasma-go/pkg/util"
	scriptbuilder "github.com/phantasma-io/phantasma-go/pkg/vm/script_builder"
)

// Put WIF and recepient address into predefinedWif and predefinedRecepient variables to skip manual console input
var predefinedWif string = ""
var predefinedRecepient string = ""

var netSelected string
var keyPair crypto.PhantasmaKeys
var client rpc.PhantasmaRPC
var chainTokens []response.TokenResult

func printTokens() {
	for _, t := range chainTokens {
		fmt.Println(t.Symbol, "flags:", t.Flags)
	}
}

func getChainToken(symbol string) response.TokenResult {
	for _, t := range chainTokens {
		if t.Symbol == symbol {
			return t
		}
	}

	panic("Token not found")
}

func menu() {
	logout := false
	for !logout {
		menuIndex, _ := PromptIndexedMenu("\nPHANTASMA GO CONSOLE WALLET DEMO. MENU:",
			[]string{"show address",
				"show balance",
				"show balance of other address",
				"send tokens",
				"staking",
				"list last 10 transactions",
				"logout"})

		switch menuIndex {
		case 1:
			showAddress()
		case 2:
			showBalance(keyPair.Address().String())
		case 3:
			showBalance(PromptStringInput("Enter address: "))
		case 4:
			sendFungibleTokens()
		case 5:
			staking()
		case 6:
			listLast10Transactions()
		case 7:
			logout = true
		}
	}
}

func showAddress() {
	fmt.Println(keyPair.Address().String())
}

func showBalance(address string) (int, []response.BalanceResult) {
	// Calling "GetAccount" method to get token balances of the address
	account, err := client.GetAccount(address)
	if err != nil {
		panic("GetAccount call failed! Error: " + err.Error())
	} else {
		fmt.Println("Balances:")

		fmt.Println("- Fungible tokens:")
		j := 1
		for i := 0; i < len(account.Balances); i += 1 {
			t := getChainToken(account.Balances[i].Symbol)
			if t.IsFungible() {
				if account.Balances[i].Symbol == "SOUL" {
					unstakedSoul := account.Balances[i].ConvertDecimalsToFloat()
					stakedSoul := account.Stakes.ConvertDecimalsToFloat()
					fmt.Printf("#%02d: %s balance: %s [not staked: %s | staked: %s]\n", j, account.Balances[i].Symbol, (new(big.Float).Add(unstakedSoul, stakedSoul)).String(), unstakedSoul.String(), stakedSoul.String())
				} else {
					fmt.Printf("#%02d: %s balance: %s\n", j, account.Balances[i].Symbol, account.Balances[i].ConvertDecimals())
				}
				j += 1
			}
		}

		fmt.Println("- Non-fungible tokens (NFTs):")
		for i := 0; i < len(account.Balances); i += 1 {
			t := getChainToken(account.Balances[i].Symbol)
			if !t.IsFungible() {
				fmt.Printf("#%02d: %s balance: %s\n", j, account.Balances[i].Symbol, account.Balances[i].ConvertDecimals())
				j += 1
			}
		}
	}

	return len(account.Balances), account.Balances
}

func sendFungibleTokens() {
	var to string
	if predefinedRecepient == "" {
		to = PromptStringInput("Enter destination address: ")
	} else {
		fmt.Println("Predefined destination address: ", predefinedRecepient)
		to = predefinedRecepient
	}

	fmt.Println("Available tokens: ")
	tokensCount, balances := showBalance(keyPair.Address().String())
	if tokensCount == 0 {
		fmt.Println("No tokens available for this address")
		return
	}
	tokenIndex := PromptIntInput("Choose token to send", 1, tokensCount)
	tokenIndex -= 1

	tokenSymbol := balances[tokenIndex].Symbol

	_, tokenAmountStr := PromptBigFloatInput("Enter amount:", big.NewFloat(0), balances[tokenIndex].ConvertDecimalsToFloat())
	tokenAmount := util.ConvertDecimalsBack(tokenAmountStr, int(balances[tokenIndex].Decimals))

	// build script
	sb := scriptbuilder.BeginScript()
	script := sb.AllowGas(keyPair.Address().String(), crypto.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
		TransferTokens(tokenSymbol, keyPair.Address().String(), to, tokenAmount).
		SpendGas(keyPair.Address().String()).
		EndScript()

	// build tx
	expire := time.Now().UTC().Add(time.Second * time.Duration(30)).Unix()
	tx := chain.NewTransaction(netSelected, "main", script, uint32(expire), []byte("GO-SDK-v0.2"))

	// sign tx
	tx.Sign(keyPair)

	fmt.Println("Tx script: " + hex.EncodeToString(script))

	// encode tx as hex
	txHex := hex.EncodeToString(tx.Bytes(true))

	fmt.Println("Tx: " + txHex)

	if !PromptYNChoice("Send transaction?") {
		return
	}

	txHash, err := client.SendRawTransaction(txHex)
	if err != nil {
		panic("Broadcasting tx failed! Error: " + err.Error())
	} else {
		if util.ErrorDetect(txHash) {
			panic("Broadcasting tx failed! Error: " + txHash)
		} else {
			fmt.Println("Tx successfully broadcasted! Tx hash: " + txHash)
		}
	}

	for {
		txResult, _ := client.GetTransaction(txHash)
		//if err != nil {
		//	fmt.Println("err: " + err.Error())
		//}
		fmt.Println("Tx state: " + fmt.Sprint(txResult.State))

		if txResult.StateIsSuccess() {
			fmt.Println("Transaction was successfully minted, tx hash: " + fmt.Sprint(txResult.Hash))
			break // Funds were transferred successfully
		}
		if txResult.StateIsFault() {
			fmt.Println("Transaction failed, tx hash: " + fmt.Sprint(txResult.Hash))
			break // Funds were not transferred
		}
	}
}

func getSoulBalanceForAddress(address string) (*big.Float, *big.Float) {
	// Calling "GetAccount" method to get token balances of the address
	account, err := client.GetAccount(address)
	if err != nil {
		panic("GetAccount call failed! Error: " + err.Error())
	}

	for i := 0; i < len(account.Balances); i += 1 {
		if account.Balances[i].Symbol == "SOUL" {
			return account.Balances[i].ConvertDecimalsToFloat(), account.Stakes.ConvertDecimalsToFloat()
		}
	}

	return big.NewFloat(0), big.NewFloat(0)
}

func staking() {
	unstakedSoul, stakedSoul := getSoulBalanceForAddress(keyPair.Address().String())
	fmt.Printf("SOUL balance: %s [not staked: %s | staked: %s]\n", (new(big.Float).Add(unstakedSoul, stakedSoul)).String(), unstakedSoul.String(), stakedSoul.String())

	menuIndex, _ := PromptIndexedMenu("SOUL STAKING MENU:", []string{"stake", "unstake", "go back"})

	stakeMode := true
	switch menuIndex {
	case 1:
		stakeMode = true
	case 2:
		stakeMode = false
	case 3:
		return
	}

	t := getChainToken("SOUL")

	var amountLimit *big.Float
	if stakeMode {
		amountLimit = unstakedSoul
	} else {
		amountLimit = stakedSoul
	}
	_, tokenAmountStr := PromptBigFloatInput("Enter amount:", big.NewFloat(0), amountLimit)
	tokenAmount := util.ConvertDecimalsBack(tokenAmountStr, int(t.Decimals))

	// build script
	sb := scriptbuilder.BeginScript()
	sb = sb.AllowGas(keyPair.Address().String(), crypto.NullAddress().String(), big.NewInt(100000), big.NewInt(21000))
	if stakeMode {
		sb = sb.Stake(keyPair.Address().String(), tokenAmount)
	} else {
		sb = sb.Unstake(keyPair.Address().String(), tokenAmount)
	}

	sb = sb.SpendGas(keyPair.Address().String())
	script := sb.EndScript()

	// build tx
	expire := time.Now().UTC().Add(time.Second * time.Duration(30)).Unix()
	tx := chain.NewTransaction(netSelected, "main", script, uint32(expire), []byte("GO-SDK-v0.2"))

	// sign tx
	tx.Sign(keyPair)

	fmt.Println("Tx script: " + hex.EncodeToString(script))

	// encode tx as hex
	txHex := hex.EncodeToString(tx.Bytes(true))

	fmt.Println("Tx: " + txHex)

	if !PromptYNChoice("Send transaction?") {
		return
	}

	txHash, err := client.SendRawTransaction(txHex)
	if err != nil {
		panic("Broadcasting tx failed! Error: " + err.Error())
	} else {
		if util.ErrorDetect(txHash) {
			panic("Broadcasting tx failed! Error: " + txHash)
		} else {
			fmt.Println("Tx successfully broadcasted! Tx hash: " + txHash)
		}
	}

	for {
		txResult, _ := client.GetTransaction(txHash)
		//if err != nil {
		//	fmt.Println("err: " + err.Error())
		//}
		fmt.Println("Tx state: " + fmt.Sprint(txResult.State))

		if txResult.StateIsSuccess() {
			fmt.Println("Transaction was successfully minted, tx hash: " + fmt.Sprint(txResult.Hash))
			break // Funds were transferred successfully
		}
		if txResult.StateIsFault() {
			fmt.Println("Transaction failed, tx hash: " + fmt.Sprint(txResult.Hash))
			break // Funds were not transferred
		}
	}
}

func listLast10Transactions() {
	// Calling "GetAddressTransactionCount" method to get transactions for the address
	transactionCount, err := client.GetAddressTransactionCount(keyPair.Address().String(), "main")
	if err != nil {
		panic("GetAddressTransactionCount call failed! Error: " + err.Error())
	} else {
		fmt.Println("Transactions count: ", transactionCount)
	}

	// Calling "GetAddressTransactions" method to get transactions for the address
	transactions, err := client.GetAddressTransactions(keyPair.Address().String(), 1, 10)
	if err != nil {
		panic("GetAddressTransactions call failed! Error: " + err.Error())
	} else {
		fmt.Println("Last 10 transactions:")
		txs := transactions.Result.(*response.AddressTransactionsResult).Txs
		for i := 0; i < len(txs); i += 1 {
			fmt.Println("#", i+1, ": ", txs[i].Hash, " timestamp: ", time.Unix(int64(txs[i].Timestamp), 0))
		}
	}
}

func main() {
	_, netSelected = PromptIndexedMenu("SELECT TESTNET OR MAINNET", []string{"testnet", "mainnet"})

	if netSelected == "testnet" {
		client = rpc.NewRPCTestnet()
	} else {
		client = rpc.NewRPCMainnet()
	}

	var err error
	chainTokens, err = client.GetTokens(false)
	fmt.Println("Received information about", len(chainTokens), netSelected, "tokens")

	// printTokens()
	// t := getChainToken("SOUL")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())
	// t = getChainToken("CROWN")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())
	// t = getChainToken("KCAL")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())

	var wif string
	if predefinedWif == "" {
		wif = PromptStringInput("Enter your WIF: ")
	} else {
		fmt.Println("Predefined WIF: ", predefinedWif)
		wif = predefinedWif
	}

	// create key pair from WIF
	keyPair, err = crypto.FromWIF(wif)
	if err != nil {
		panic("Creating keyPair failed!")
	}

	menu()
}
