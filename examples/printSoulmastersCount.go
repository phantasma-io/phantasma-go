package main

import (
	"encoding/hex"
	"fmt"

	"github.com/phantasma-io/phantasma-go/pkg/io"
	"github.com/phantasma-io/phantasma-go/pkg/vm"
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

	decoded, _ := hex.DecodeString(result.Result)
	br := io.NewBinReaderFromBuf(decoded)

	var vmObject vm.VMObject
	vmObjectResult := vmObject.Deserialize(br)

	if err != nil {
		panic("Script invocation failed! Error: " + err.Error())
	}

	fmt.Println("Script invocation result: ", vmObjectResult.AsNumber().String())
}
