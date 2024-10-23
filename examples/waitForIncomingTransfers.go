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
	// Get current block height
	height, _ := client.GetBlockHeight("main")

	for {
		fmt.Println("Checking new block #", height.String())
		// Get block's data by its height
		block, err := client.GetBlockByHeight("main", height.String())
		if err != nil {
			panic("GetBlockByHeight call failed! Error: " + err.Error())
		}

		// Iterate throough all transactions in the block
		for _, tx := range block.Txs {
			// Skip failed trasactions
			if !tx.StateIsSuccess() {
				continue
			}

			// Iterate throough all events in the transaction
			for _, e := range tx.Events {

				if e.Kind == event.TokenReceive.String() && e.Address == address {
					// We found TokenReceive event for given address

					// Decode event data into event.TokenEventData structure
					decoded, _ := hex.DecodeString(e.Data)
					data := io.Deserialize[*event.TokenEventData](decoded)

					// Apply decimals to the token amount
					t := getChainToken(data.Symbol)
					tokenAmount := util.ConvertDecimals(data.Value, int(t.Decimals))

					// Call our callback function
					onTransactionReceived(e.Address, data.Symbol, tokenAmount)
				}
			}
		}

		// Wait for next block to appear on the blockchain
		for {
			newHeight, _ := client.GetBlockHeight("main")
			if newHeight.Cmp(height) == 1 {
				// New block was minted (at least 1 new block)
				height = height.Add(height, big.NewInt(1))
				break
			}

			// Wait 200 milliseconds before making next RPC call
			time.Sleep(200 * time.Millisecond)
		}
	}
}
