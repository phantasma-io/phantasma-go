package main

import (
	"math/big"
)

func getSoulBalance(address string) (*big.Float, *big.Float) {
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
