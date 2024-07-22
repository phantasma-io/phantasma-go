package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
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
	reader := bufio.NewReader(os.Stdin)

	logout := false
	for !logout {
		fmt.Println()
		fmt.Println("PHANTASMA GO CONSOLE WALLET DEMO. MENU:")
		fmt.Println("1 - show address")
		fmt.Println("2 - show balance")
		fmt.Println("3 - show balance of other address")
		fmt.Println("4 - send tokens")
		fmt.Println("5 - staking")
		fmt.Println("6 - list last 10 transactions")
		fmt.Println("7 - logout")

		menuIndexStr, _ := reader.ReadString('\n')
		menuIndexStr = strings.TrimSuffix(menuIndexStr, "\n")
		menuIndex, _ := strconv.Atoi(menuIndexStr)

		switch menuIndex {
		case 1:
			showAddress()
		case 2:
			showBalance(keyPair.Address().String())
		case 3:
			fmt.Print("Enter address: ")
			address, _ := reader.ReadString('\n')
			address = strings.TrimSuffix(address, "\n")
			showBalance(address)
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
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter destination address: ")
	var to string
	if predefinedRecepient == "" {
		to, _ = reader.ReadString('\n')
		to = strings.TrimSuffix(to, "\n")
	} else {
		fmt.Println(predefinedRecepient)
		to = predefinedRecepient
	}

	fmt.Println("Available tokens: ")
	tokensCount, balances := showBalance(keyPair.Address().String())
	if tokensCount == 0 {
		fmt.Println("No tokens available for this address")
		return
	}
	fmt.Print("Choose token to send (#1-", tokensCount, "): ")
	tokenIndexStr, _ := reader.ReadString('\n')
	tokenIndexStr = strings.TrimSuffix(tokenIndexStr, "\n")
	tokenIndex, _ := strconv.Atoi(tokenIndexStr)
	if tokenIndex > tokensCount {
		fmt.Println("Incorrect token number entered")
		return
	}
	tokenIndex -= 1

	tokenSymbol := balances[tokenIndex].Symbol

	fmt.Print("Enter amount: (max ", balances[tokenIndex].ConvertDecimals(), "): ")
	tokenAmountStr, _ := reader.ReadString('\n')
	tokenAmountStr = strings.TrimSuffix(tokenAmountStr, "\n")
	tokenAmount := util.ConvertDecimalsBack(tokenAmountStr, int(balances[tokenIndex].Decimals))

	// build script
	sb := scriptbuilder.BeginScript()
	script := sb.AllowGas(keyPair.Address().String(), crypto.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
		TransferTokens(tokenSymbol, keyPair.Address().String(), to, tokenAmount).
		SpendGas(keyPair.Address().String()).
		EndScript()

	// build tx
	expire := time.Now().UTC().Add(time.Second * time.Duration(30)).Unix()
	tx := chain.NewTransaction("mainnet", "main", script, uint32(expire), []byte("GO-SDK-v0.2"))

	// sign tx
	tx.Sign(keyPair)

	fmt.Println("Tx script: " + hex.EncodeToString(script))

	// encode tx as hex
	txHex := hex.EncodeToString(tx.Bytes(true))

	fmt.Println("Tx: " + txHex)

	for true {
		fmt.Print("Send transaction? (y/n): ")
		sendTransactionYN, _ := reader.ReadString('\n')
		sendTransactionYN = strings.TrimSuffix(sendTransactionYN, "\n")
		if strings.ToLower(sendTransactionYN) == "n" {
			return
		}
		if strings.ToLower(sendTransactionYN) == "y" {
			break
		}
		fmt.Println("Please enter 'y' or 'n'")
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
	reader := bufio.NewReader(os.Stdin)

	unstakedSoul, stakedSoul := getSoulBalanceForAddress(keyPair.Address().String())
	fmt.Printf("SOUL balance: %s [not staked: %s | staked: %s]\n", (new(big.Float).Add(unstakedSoul, stakedSoul)).String(), unstakedSoul.String(), stakedSoul.String())

	fmt.Println("SOUL STAKING MENU:")
	fmt.Println("1 - stake")
	fmt.Println("2 - unstake")
	fmt.Println("3 - go back")

	menuIndexStr, _ := reader.ReadString('\n')
	menuIndexStr = strings.TrimSuffix(menuIndexStr, "\n")
	menuIndex, _ := strconv.Atoi(menuIndexStr)

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

	var tokenAmount *big.Int
	for true {
		if stakeMode {
			fmt.Print("Enter amount: (max ", unstakedSoul, "): ")
		} else {
			fmt.Print("Enter amount: (max ", stakedSoul, "): ")
		}
		tokenAmountStr, _ := reader.ReadString('\n')
		tokenAmountStr = strings.TrimSuffix(tokenAmountStr, "\n")
		tokenAmount = util.ConvertDecimalsBack(tokenAmountStr, int(t.Decimals))

		tokenAmountFloat, _ := big.NewFloat(0).SetString(tokenAmountStr)
		if stakeMode {
			if tokenAmountFloat.Cmp(unstakedSoul) > 0 {
				fmt.Printf("Entered amount '%s' is higher than SOUL amount available for staking '%s'\n", tokenAmountFloat.String(), unstakedSoul.String())
			} else {
				break
			}
		} else {
			if tokenAmountFloat.Cmp(stakedSoul) > 0 {
				fmt.Printf("Entered amount '%s' is higher than SOUL amount available for unstaking '%s'\n", tokenAmountFloat.String(), stakedSoul.String())
			} else {
				break
			}
		}
	}

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
	tx := chain.NewTransaction("mainnet", "main", script, uint32(expire), []byte("GO-SDK-v0.2"))

	// sign tx
	tx.Sign(keyPair)

	fmt.Println("Tx script: " + hex.EncodeToString(script))

	// encode tx as hex
	txHex := hex.EncodeToString(tx.Bytes(true))

	fmt.Println("Tx: " + txHex)

	for true {
		fmt.Print("Send transaction? (y/n): ")
		sendTransactionYN, _ := reader.ReadString('\n')
		sendTransactionYN = strings.TrimSuffix(sendTransactionYN, "\n")
		if strings.ToLower(sendTransactionYN) == "n" {
			return
		}
		if strings.ToLower(sendTransactionYN) == "y" {
			break
		}
		fmt.Println("Please enter 'y' or 'n'")
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
	client = rpc.NewRPCMainnet()

	var err error
	chainTokens, err = client.GetTokens(false)
	fmt.Println("Received information about", len(chainTokens), "chain tokens")

	// printTokens()
	// t := getChainToken("SOUL")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())
	// t = getChainToken("CROWN")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())
	// t = getChainToken("KCAL")
	// fmt.Println(t.Symbol, "fungible:", t.IsFungible(), "fuel:", t.IsFuel(), "stakable:", t.IsStakable(), "burnable:", t.IsBurnable(), "transferable:", t.IsTransferable())

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your WIF: ")
	var wif string
	if predefinedWif == "" {
		wif, _ = reader.ReadString('\n')
		wif = strings.TrimSuffix(wif, "\n")
	} else {
		fmt.Println(predefinedWif)
		wif = predefinedWif
	}

	// create key pair from WIF
	keyPair, err = crypto.FromWIF(wif)
	if err != nil {
		panic("Creating keyPair failed!")
	}

	menu()
}
