package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/phantasma-io/phantasma-go/pkg/domain/event"
	"github.com/phantasma-io/phantasma-go/pkg/io"
	"github.com/phantasma-io/phantasma-go/pkg/util"
)

func onTransactionReceived(address, symbol, amount string) {
	fmt.Printf("Address %s received %s %s\n", address, amount, symbol)
}

func waitForIncomingTransfers(address string) {
	height, _ := client.GetBlockHeight("main")

	for {
		fmt.Println("Checking new block #", height.String())
		block, err := client.GetBlockByHeight("main", height.String())
		if err != nil {
			panic("GetBlockByHeight call failed! Error: " + err.Error())
		}

		for _, tx := range block.Txs {
			for _, e := range tx.Events {

				if e.Kind == event.TokenReceive.String() && e.Address == address {
					decoded, _ := hex.DecodeString(e.Data)
					br := io.NewBinReaderFromBuf(decoded)

					var data event.TokenEventData
					data.Deserialize(br)

					t := getChainToken(data.Symbol)
					tokenAmount := util.ConvertDecimals(data.Value, int(t.Decimals))
					onTransactionReceived(e.Address, data.Symbol, tokenAmount)
				}
			}
		}

		for {
			newHeight, _ := client.GetBlockHeight("main")
			if newHeight.Cmp(height) == 1 {
				height = height.Add(height, big.NewInt(1))
				break
			}

			time.Sleep(200 * time.Millisecond)
		}
	}
}
