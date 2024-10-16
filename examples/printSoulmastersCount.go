package main

import (
	"encoding/hex"
	"fmt"

	scriptbuilder "github.com/phantasma-io/phantasma-go/pkg/vm/script_builder"
)

func printSoulmastersCount() {
	// Build script
	sb := scriptbuilder.BeginScript().
		CallContract("stake", "GetMasterCount")
	script := sb.EndScript()

	// Before sending script to the chain we need to encode it into Base16 encoding (HEX)
	encodedScript := hex.EncodeToString(script)
	fmt.Println("Script: " + encodedScript)

	if !PromptYNChoice("Invoke script?") {
		return
	}

	// Make the call itself
	result, err := client.InvokeRawScript("main", encodedScript)

	if err != nil {
		panic("Script invocation failed! Error: " + err.Error())
	}

	fmt.Println("Current SoulMasters count: ", result.DecodeResult().AsNumber().String())
}

func printSoulmastersCountAndLastInflationDate() {
	// Build script
	sb := scriptbuilder.BeginScript().
		CallContract("stake", "GetMasterCount")

	sb.CallContract("gas", "GetLastInflationDate")

	script := sb.EndScript()

	// Before sending script to the chain we need to encode it into Base16 encoding (HEX)
	encodedScript := hex.EncodeToString(script)
	fmt.Println("Script: " + encodedScript)

	if !PromptYNChoice("Invoke script?") {
		return
	}

	// Make the call itself
	result, err := client.InvokeRawScript("main", encodedScript)

	if err != nil {
		panic("Script invocation failed! Error: " + err.Error())
	}

	fmt.Printf("Current SoulMasters count: %s, last inflation date: %s \n", result.DecodeResults(0).AsNumber().String(), result.DecodeResults(1).AsNumber().String())
}
