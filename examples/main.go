package main

import (
	"fmt"
	"math/big"

	crypto "github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/rpc"
	"github.com/phantasma-io/phantasma-go/pkg/rpc/response"
	"github.com/phantasma-io/phantasma-go/pkg/util"
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
			printBalance(keyPair.Address().String())
		case 3:
			printBalance(PromptStringInput("Enter address: "))
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

func sendFungibleTokens() {
	var to string
	if predefinedRecepient == "" {
		to = PromptStringInput("Enter destination address: ")
	} else {
		fmt.Println("Predefined destination address: ", predefinedRecepient)
		to = predefinedRecepient
	}

	fmt.Println("Available tokens: ")
	tokensCount, balances := printBalance(keyPair.Address().String())
	if tokensCount == 0 {
		fmt.Println("No tokens available for this address")
		return
	}
	tokenIndex := PromptIntInput("Choose token to send", 1, tokensCount)
	tokenIndex -= 1

	tokenSymbol := balances[tokenIndex].Symbol

	_, tokenAmountStr := PromptBigFloatInput("Enter amount:", big.NewFloat(0), balances[tokenIndex].ConvertDecimalsToFloat())
	tokenAmount := util.ConvertDecimalsBack(tokenAmountStr, int(balances[tokenIndex].Decimals))

	sendFungibleToken(tokenSymbol, to, tokenAmount)
}

func staking() {
	unstakedSoul, stakedSoul := getSoulBalance(keyPair.Address().String())
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

	if stakeMode {
		stakeSoulToken(keyPair.Address().String(), tokenAmount)
	} else {
		unstakeSoulToken(keyPair.Address().String(), tokenAmount)
	}
}

func listLast10Transactions() {
	printTransactionsCount(keyPair.Address().String())
	printTransactions(keyPair.Address().String(), 1, 10)
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
