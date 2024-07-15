package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	chain "github.com/phantasma-io/phantasma-go/pkg/blockchain"
	crypto "github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/rpc"
	"github.com/phantasma-io/phantasma-go/pkg/util"
	scriptbuilder "github.com/phantasma-io/phantasma-go/pkg/vm/script_builder"
)

func main() {

	// create key pair from WIF
	kp, err := crypto.FromWIF("ADD_WIF_HERE")
	if err != nil {
		panic("Creating keyPair failed!")
	}

	// build script
	sb := scriptbuilder.BeginScript()
	script := sb.AllowGas(kp.Address(), crypto.NullAddress(), *big.NewInt(100000), *big.NewInt(21000)).
		TransferTokens("SOUL", kp.Address().String(), "ADD_RECEIVER_HERE", *big.NewInt(100000000)).
		SpendGas(kp.Address()).
		EndScript()

	// build tx
	expire := time.Now().UTC().Add(time.Second * time.Duration(30)).Unix()
	tx := chain.NewTransaction("mainnet", "main", script, uint32(expire), []byte("GO-SDK-v0.2"))

	// sign tx
	tx.Sign(kp)

	fmt.Println("Tx script: " + hex.EncodeToString(script))

	// encode tx as hex
	txHex := hex.EncodeToString(tx.Bytes(true))

	fmt.Println("Tx: " + txHex)

	client := rpc.NewRPCMainnet()
	txHash, err := client.SendRawTransaction(txHex)
	if err != nil {
		panic("Broadcasting tx failed! Error: " + err.Error())
	} else {
		if util.ErrorDetect(txHash) {
			panic("Broadcasting tx failed! Error: " + txHash)
		} else {
			fmt.Println("Tx successfully broadcasted!")
		}
	}

	for {
		txResult, _ := client.GetTransaction(txHash)
		//if err != nil {
		//	fmt.Println("err: " + err.Error())
		//}
		fmt.Println("Tx hash: " + fmt.Sprint(txResult.Hash))

		if txResult.Hash != "" {
			fmt.Println("Transaction was successfully minted, tx hash: " + fmt.Sprint(txResult.Hash))
			break // Funds were transferred successfully
		}
	}
}
