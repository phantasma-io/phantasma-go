package main

import (
	"fmt"

	"github.com/phantasma-io/phantasma-go/pkg/rpc"
	"github.com/phantasma-io/phantasma-go/pkg/rpc/response"
)

var netSelected string
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
		menuIndex, _ := PromptIndexedMenu("\nPHANTASMA GO CONSOLE DEMO. MENU:",
			[]string{"wallet",
				"show balance of address",
				"chain stats",
				"logout"})

		switch menuIndex {
		case 1:
			wallet()
		case 2:
			printBalance(PromptStringInput("Enter address: "))
		case 3:
			chainStats()
		case 4:
			logout = true
		}
	}
}

func chainStats() {
	menuIndex, _ := PromptIndexedMenu("CHAIN STATS MENU:", []string{"latest block height", "soulmasters count", "soulmasters count and last inflation date", "go back"})

	switch menuIndex {
	case 1:
		height, _ := client.GetBlockHeight("main")
		fmt.Println("Latest block height:", height)
	case 2:
		printSoulmastersCount()
	case 3:
		printSoulmastersCountAndLastInflationDate()
	case 4:
		return
	}
}

func main() {
	_, netSelected = PromptIndexedMenu("SELECT TESTNET OR MAINNET", []string{"testnet", "mainnet"})

	if netSelected == "testnet" {
		client = rpc.NewRPCTestnet()
	} else {
		client = rpc.NewRPCMainnet()
	}

	chainTokens, _ = client.GetTokens(false)
	fmt.Println("Received information about", len(chainTokens), netSelected, "tokens")

	// printTokens()
	// t := getChainToken("SOUL")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())
	// t = getChainToken("CROWN")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())
	// t = getChainToken("KCAL")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())

	menu()
}
