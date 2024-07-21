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
		fmt.Println("5 - list last 10 transactions")
		fmt.Println("6 - logout")

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
			listLast10Transactions()
		case 6:
			logout = true
			break
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
		for i := 0; i < len(account.Balances); i += 1 {
			fmt.Println("#", i+1, ": ", account.Balances[i].Symbol, " balance: ", account.Balances[i].ConvertDecimals())
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
	var err error
	keyPair, err = crypto.FromWIF(wif)
	if err != nil {
		panic("Creating keyPair failed!")
	}

	client = rpc.NewRPCMainnet()

	menu()
}
