package main

import (
	"encoding/hex"
	"fmt"

	scriptbuilder "github.com/phantasma-io/phantasma-go/pkg/vm/script_builder"
)

func printSoulmastersCount() {
	// build script
	sb := scriptbuilder.BeginScript().
		CallContract("stake", "GetMasterCount")
	script := sb.EndScript()

	encodedScript := hex.EncodeToString(script)
	fmt.Println("Script: " + encodedScript)

	if !PromptYNChoice("Invoke script?") {
		return
	}

	result, err := client.InvokeRawScript("main", encodedScript)

	if err != nil {
		panic("Script invocation failed! Error: " + err.Error())
	}

	fmt.Println("Current SoulMasters count: ", result.DecodeResult().AsNumber().String())
}
