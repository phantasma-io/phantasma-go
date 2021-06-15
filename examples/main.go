package main

import (
	"encoding/hex"
	"fmt"
	"time"

	chain "github.com/phantasma.io/phantasma-go/pkg/blockchain"
	crypto "github.com/phantasma.io/phantasma-go/pkg/cryptography"
	"github.com/phantasma.io/phantasma-go/pkg/rpc"
	scriptbuilder "github.com/phantasma.io/phantasma-go/pkg/vm/script_builder"
)

func main() {

	// create key pair from WIF
	kp, err := crypto.FromWIF("ADD_WIF_HERE")
	if err != nil {
		panic("Creating keyPair failed!")
	}

	// build script
	sb := scriptbuilder.BeginScript()
	script := sb.AllowGas(kp.Address().String(), "", 100000, 1200).
		TransferTokens("SOUL", kp.Address().String(), "ADD_RECEIVER_HERE", 100000000).
		SpendGas(kp.Address().String()).
		EndScript()

	// build tx
	expire := time.Now().Add(time.Second * time.Duration(30)).Unix()
	tx := chain.NewTransaction("mainnet", "main", script, uint32(expire), []byte("GO-SDK-v0.1"))

	// sign tx
	tx.Sign(kp)

	// encode tx as hex
	txHex := hex.EncodeToString(tx.Bytes(true))
	client := rpc.NewRPCMainnet()
	txHash, err := client.SendRawTransaction(txHex)
	if err != nil {
		panic("Broadcasting tx " + txHash + " failed!")
	}

	fmt.Println("tx " + txHash + " successfully broadcasted!")

	for {
		txResult, _ := client.GetTransaction(txHash)
		//if err != nil {
		//	fmt.Println("err: " + err.Error())
		//}
		fmt.Println("txHeight: " + fmt.Sprint(txResult.Hash))
	}
}
