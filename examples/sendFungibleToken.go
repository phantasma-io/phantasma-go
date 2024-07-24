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

func sendFungibleToken(tokenSymbol, to string, tokenAmount *big.Int) {
	// build script
	sb := scriptbuilder.BeginScript()
	script := sb.AllowGas(keyPair.Address().String(), crypto.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
		TransferTokens(tokenSymbol, keyPair.Address().String(), to, tokenAmount).
		SpendGas(keyPair.Address().String()).
		EndScript()

	// build tx
	expire := time.Now().UTC().Add(time.Second * time.Duration(30)).Unix()
	tx := chain.NewTransaction(netSelected, "main", script, uint32(expire), domain.SDKPayload)

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

	waitForTransactionResult(txHash)
}
