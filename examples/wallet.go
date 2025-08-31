package main

import (
	"fmt"
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/util"
)

// Put WIF and recepient address into predefinedWif and predefinedRecepient variables to skip manual console input
var predefinedWif string = ""
var predefinedRecepient string = ""

var keyPair cryptography.PhantasmaKeys

func wallet() {
	var wif string
	if predefinedWif == "" {
		wif = PromptStringInput("Enter your WIF: ")
	} else {
		fmt.Println("Predefined WIF: ", predefinedWif)
		wif = predefinedWif
	}

	// create key pair from WIF
	var err error
	keyPair, err = cryptography.FromWIF(wif)
	if err != nil {
		panic("Creating keyPair failed!")
	}

	goBack := false
	for !goBack {
		menuIndex, _ := PromptIndexedMenu("\nPHANTASMA WALLET DEMO. MENU:",
			[]string{"show address",
				"show balance",
				"wait for incoming transfers",
				"send tokens",
				"staking",
				"list last 10 transactions",
				"go back"})

		switch menuIndex {
		case 1:
			showAddress()
		case 2:
			printBalance(keyPair.Address().String())
		case 3:
			waitForIncomingTransfers(keyPair.Address().String())
		case 4:
			sendFungibleTokens()
		case 5:
			staking()
		case 6:
			listLast10Transactions()
		case 7:
			goBack = true
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

	toAddress, _ := cryptography.FromString(to)
	sendFungibleToken(tokenSymbol, toAddress, tokenAmount)
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
		stakeSoulToken(keyPair.Address(), tokenAmount)
	} else {
		unstakeSoulToken(keyPair.Address(), tokenAmount)
	}
}

func listLast10Transactions() {
	printTransactionsCount(keyPair.Address().String())
	printTransactions(keyPair.Address().String(), 1, 10)
}
