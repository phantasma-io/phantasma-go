package main

import (
	"fmt"
	"time"
)

func printTransactionsCount(address string) {
	// Calling "GetAddressTransactionCount" method to get transactions for the address
	transactionCount, err := client.GetAddressTransactionCount(address, "main")
	if err != nil {
		panic("GetAddressTransactionCount call failed! Error: " + err.Error())
	} else {
		fmt.Println("Transactions count: ", transactionCount)
	}
}

func printTransactions(address string, page, pageSize int) {
	// Calling "GetAddressTransactions" method to get transactions for the address
	transactions, err := client.GetAddressTransactions(address, page, pageSize)
	if err != nil {
		panic("GetAddressTransactions call failed! Error: " + err.Error())
	} else {
		fmt.Println("Transactions:")
		txs := transactions.Result.Txs
		for i := 0; i < len(txs); i += 1 {
			fmt.Println("#", i+1, ": ", txs[i].Hash, " timestamp: ", time.Unix(int64(txs[i].Timestamp), 0))
		}
	}
}
