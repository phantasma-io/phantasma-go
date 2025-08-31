package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	chain "github.com/phantasma-io/phantasma-go/pkg/blockchain"
	crypto "github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/domain"
	"github.com/phantasma-io/phantasma-go/pkg/util"
	scriptbuilder "github.com/phantasma-io/phantasma-go/pkg/vm/script_builder"
)

func stakeSoulToken(address crypto.Address, tokenAmount *big.Int) {
	// Build script
	sb := scriptbuilder.BeginScript().
		AllowGas(address, crypto.NullAddress(), big.NewInt(100000), big.NewInt(21000)).
		Stake(address, tokenAmount).
		SpendGas(address)
	script := sb.EndScript()

	// Build transaction
	expire := time.Now().UTC().Add(time.Second * time.Duration(30)).Unix()
	tx := chain.NewTransaction(netSelected, "main", script, uint32(expire), domain.SDKPayload)

	// Sign transaction
	tx.Sign(keyPair)

	fmt.Println("Tx script: " + hex.EncodeToString(script))

	// Before sending script to the chain we need to encode it into Base16 encoding (HEX)
	txHex := hex.EncodeToString(tx.Bytes())

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

	waitForTransactionResult(txHash)
}
