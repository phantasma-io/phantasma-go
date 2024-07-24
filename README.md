<p align="center">
<img src="./.github/phantasma-go.jpg" width="300px" alt="logo">
</p>
<p align="center">
  <b>Go</b> SDK for the <a href="https://phantasma.io">Phantasma</a> blockchain.
</p>

<hr />

![License](https://img.shields.io/github/license/phantasma-io/phantasma-go.svg?style=popout)

# Overview

This project aims to be an easy to use SDK for the Phantasma blockchain.

# Documentation

## Installation

PhantasmaGo is distributed as a library that includes all the functionality provided.

```
go get -u github.com/phantasma-io/phantasma-go
```

## Getting started
To start interacting with Phantasma blockchain you need to choose network you are planning to use (mainnet or testnet) and create corresponding RPC client.

Creation of testnet RPC client:
```
client = rpc.NewRPCTestnet()
```

Creation of mainnet RPC client:
```
client = rpc.NewRPCMainnet()
```

To create a new key pair structure from private key in WIF format use following code:
```
keyPair, err := cryptography.FromWIF("put WIF here")
if err != nil {
    panic("Creating keyPair failed!")
}
```

To get detailed description of tokens deployed on the chain you can use following code:

```
var chainTokens []response.TokenResult

func getChainToken(symbol string) response.TokenResult {
    for _, t := range chainTokens {
        if t.Symbol == symbol {
            return t
        }
    }

    panic("Token not found")
}

chainTokens, _ = client.GetTokens(false)
```

This will allow you to get token characteristics this way:

```
t := getChainToken("SOUL")
if t.IsFungible() {
    fmt.Println("Token SOUL is fungible")
}
```

Code samples in the following sections of this documentation use `client` and `keyPair` structures and method `getChainToken` which should be initialized in advance.

## Script Builder

Building a script is the most important part of interacting with the Phantasma blockchain. Without a propper script, the Phantasma blockchain will not know what you are trying to do.

These functions, `CallContract` and `CallInterop`, are your bread and butter for creating new scripts.

```
func (s ScriptBuilder) CallContract(contractName, method string, args ...interface{})
```

```
func (s ScriptBuilder) CallInterop(method string, args ...interface{})
```

You can find out all the diffrent `CallInterop` functions below.

For `CallContract`, you will have to look through the ABI's of all the diffrent smart contracts currently deployed on the Phantasma 'mainnet': [Link Here](https://explorer.phantasma.info/en/nexus?tab=contracts). To see all methods of a contract, for example `stake`, you can check it with explorer: [Link Here](https://explorer.phantasma.info/en/contract?id=stake&tab=methods).

### Examples

Following code generates script to transfer `tokenAmount` amount of token `tokenSymbol` from wallet `from` to wallet `to`
```
from := "put sender address here" // Phantasma address, starting with capital 'P'
to := "put recepient address here" // Phantasma address, starting with capital 'P'
tokenAmount := big.NewInt(1000000000) // Token amount in the form of big integer
tokenSymbol := "SOUL"

sb := scriptbuilder.BeginScript()
script := sb.CallContract("gas", "AllowGas", from, cryptography.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
    CallInterop("Runtime.TransferTokens", from, to, tokenSymbol, tokenAmount).
    CallContract("gas", "SpendGas", from).
    EndScript()
```

And here we generate script to make a call which does not require transaction, for this we use `CallContract` method:

```
address := "put caller address here" // Phantasma address, starting with capital 'P'
tokenAmount := big.NewInt(1000000000) // Token amount in the form of big integer

sb := scriptbuilder.BeginScript().
    CallContract("gas", "AllowGas", address, cryptography.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
    CallContract("stake", "Stake", address, tokenAmount).
    CallContract("gas", "SpendGas", address)
script := sb.EndScript()
```

## Script Builder Extensions

For some widely used contract calls SDK has special extension methods which make code more compact. Here's the list of available extensions:

```
func (s ScriptBuilder) AllowGas(from, to string, gasPrice, gasLimit *big.Int)
```

```
func (s ScriptBuilder) SpendGas(address string)
```

```
func (s ScriptBuilder) MintTokens(symbol, from, to string, amount *big.Int)
```

```
func (s ScriptBuilder) Stake(address string, amount *big.Int)
```

```
func (s ScriptBuilder) Unstake(address string, amount *big.Int)
```

```
func (s ScriptBuilder) TransferTokens(symbol, from, to string, amount *big.Int)
```

```
func (s ScriptBuilder) TransferBalance(symbol, from, to string)
```

### Examples

We can rewrite examples from previous section using `AllowGas()` and `SpendGas()` extensions:

```
sb := scriptbuilder.BeginScript()
script := sb.AllowGas(from, cryptography.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
    CallInterop("Runtime.TransferTokens", from, to, tokenSymbol, tokenAmount).
    SpendGas(from).
    EndScript()
```
```
sb := scriptbuilder.BeginScript().
    AllowGas(address, crypto.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
    CallContract("stake", "Stake", address, tokenAmount).
    SpendGas(address)
script := sb.EndScript()
```

We can also rewrite main contract calls in these examples:

```
sb := scriptbuilder.BeginScript()
script := sb.AllowGas(from, cryptography.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
    TransferTokens(from, to, tokenSymbol, tokenAmount).
    SpendGas(from).
    EndScript()
```
```
sb := scriptbuilder.BeginScript().
    AllowGas(address, crypto.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
    Stake(address, tokenAmount).
    SpendGas(address)
script := sb.EndScript()
```

## InvokeRawScript and decoding the result

Scripts which does not require transaction can be sent to the chain directly using `InvokeRawScript()` call.

Here's an example of such call to get SoulMaster count from the chain:

```
// Build script
sb := scriptbuilder.BeginScript().
    CallContract("stake", "GetMasterCount")
script := sb.EndScript()

// Before sending script to the chain we need to encode it into Base16 encoding (HEX)
encodedScript := hex.EncodeToString(script)

// Make the call itself
result, err := client.InvokeRawScript("main", encodedScript)

if err != nil {
    panic("Script invocation failed! Error: " + err.Error())
}

// `DecodeResult()` decodes HEX-encoded byte array result, stored in `.Result` field, into `vm.VMObject` structure
// `AsNumber()` returns value stored in `vm.VMObject` structure, in `.Data` field, as a *big.Int number (in our case value is stored in `vm.VMObject` as big integer serialized into byte array)
fmt.Println("Current SoulMasters count: ", result.DecodeResult().AsNumber().String())
```

## Building and sending transaction

### Building transaction
To build a transaction you will first need to build a script.

Note, building a transaction is for transactional scripts only. Non transactional scripts should use the RPC function `InvokeRawScript()`.

```
// Build script
sb := scriptbuilder.BeginScript()
script := sb.AllowGas(keyPair.Address().String(), cryptography.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
    TransferTokens(tokenSymbol, keyPair.Address().String(), to, tokenAmount).
    SpendGas(keyPair.Address().String()).
    EndScript()

// Build transaction
expire := time.Now().UTC().Add(time.Second * time.Duration(30)).Unix()
tx := blockchain.NewTransaction(netSelected, "main", script, uint32(expire), domain.SDKPayload)

// Sign transaction
tx.Sign(keyPair)

// Before sending script to the chain we need to encode it into Base16 encoding (HEX)
txHex := hex.EncodeToString(tx.Bytes(true))
```

### Sending transaction

Here we send transaction prepared in previous block of code and stored as HEX in `txHex` variable.

```
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
```

### Waiting for transaction execution result

We need to wait for transaction to be minted on the chain to get its status:

```
for {
    txResult, _ := client.GetTransaction(txHash)

    if txResult.StateIsSuccess() {
        fmt.Println("Transaction was successfully minted, tx hash: " + fmt.Sprint(txResult.Hash))
        break // Funds were transferred successfully
    }
    if txResult.StateIsFault() {
        fmt.Println("Transaction failed, tx hash: " + fmt.Sprint(txResult.Hash))
        break // Funds were not transferred, transaction failed
    }

    time.Sleep(200 * time.Millisecond)
}
```

## Staking SOUL token

Following code shows how to stake SOUL token:

```
// Build script
sb := scriptbuilder.BeginScript().
    AllowGas(address, crypto.NullAddress().String(), big.NewInt(100000), big.NewInt(21000)).
    Stake(address, tokenAmount).
    SpendGas(address)
script := sb.EndScript()

// Build transaction
expire := time.Now().UTC().Add(time.Second * time.Duration(30)).Unix()
tx := chain.NewTransaction(netSelected, "main", script, uint32(expire), domain.SDKPayload)

// Sign transaction
tx.Sign(keyPair)

// Before sending script to the chain we need to encode it into Base16 encoding (HEX)
txHex := hex.EncodeToString(tx.Bytes(true))

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

    if txResult.StateIsSuccess() {
        fmt.Println("Transaction was successfully minted, tx hash: " + fmt.Sprint(txResult.Hash))
        break // Funds were transferred successfully
    }
    if txResult.StateIsFault() {
        fmt.Println("Transaction failed, tx hash: " + fmt.Sprint(txResult.Hash))
        break // Funds were not transferred, transaction failed
    }

    time.Sleep(200 * time.Millisecond)
}
```

## Scanning the blockchain for incoming transactions

In the following code we monitor the blockchain by checking all the new blocks minted on the blockchain and waiting for `TokenReceive` event for given address. This event for address means that address has received some tokens.

```
func onTransactionReceived(address, symbol, amount string) {
    fmt.Printf("Address %s received %s %s\n", address, amount, symbol)
}

func waitForIncomingTransfers(address string) {
    // Get current block height
    height, _ := client.GetBlockHeight("main")

    for {
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
                    br := io.NewBinReaderFromBuf(decoded)

                    var data event.TokenEventData
                    data.Deserialize(br)

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
```

## Examples

This repository has `examples` folder with some code which can be easily reused. Examples are grouped into a single console application.

To run this application switch to `examples` folder and run:

```
go run .
```

or

```
sh run.sh
```

Application entry point is `main()` function in `main.go` source file. Once launched it will display the following menu:

![image](https://github.com/user-attachments/assets/42c7f5a5-41ea-4dfa-9a05-45028c563be6)

Wallet submenu:

![image](https://github.com/user-attachments/assets/c6484554-3be8-415d-a8a0-1d885d0f9a37)

Chain stats submenu:

![image](https://github.com/user-attachments/assets/7a1c0fec-9b92-4a8a-8e16-cd54f42b066f)


# Contributing

Feel free to contribute to this project after reading the
[contributing guidelines](CONTRIBUTING.md).

Before starting to work on a certain topic, create an new issue first,
describing the feature/topic you are going to implement.

# Contact

- Get in contact with us on the [Phantasma Discord](https://discord.gg/JzSnmFZCcD)

# License

- Open-source [MIT](LICENSE.md)
