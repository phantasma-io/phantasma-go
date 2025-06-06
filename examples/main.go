package main

import (
	"encoding/hex"
	"fmt"

	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
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
				"misc",
				"logout"})

		switch menuIndex {
		case 1:
			wallet()
		case 2:
			printBalance(PromptStringInput("Enter address: "))
		case 3:
			chainStats()
		case 4:
			misc()
		case 5:
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

func misc() {
	menuIndex, _ := PromptIndexedMenu("MISC MENU:", []string{"address from public key", "go back"})

	switch menuIndex {
	case 1:
		publicKeyHex := PromptStringInput("Enter public key in hex: ")
		publicKey, err := hex.DecodeString(publicKeyHex)
		if err != nil {
			panic(err)
		}

		if len(publicKey) == cryptography.Length {
			// This is the only correct way, address should have 34 bytes.
			// 1 byte for type, 1 byte is reserved (must be 0) and then 32 bytes of public key.
			fmt.Println("Address: ", cryptography.NewAddress(publicKey).String())
		} else if len(publicKey) == 33 {
			publicKey = append([]byte{byte(cryptography.User)}, publicKey...)
			fmt.Println("[33 bytes] * DON'T USE THIS ADDRESS * Address type is missing, using User by default: ")
			fmt.Println(cryptography.NewAddress(publicKey).String())
		} else if len(publicKey) == 32 {
			// We use only '0x00' reserved byte for 2nd byte.
			publicKey = append([]byte{byte(cryptography.User), 0x00}, publicKey...)

			fmt.Println("[32 bytes] * DON'T USE THESE ADDRESSES * Address type is missing, using User by default: ")
			fmt.Println(cryptography.NewAddress(publicKey).String())
		}

	case 2:
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
