package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func PromptIndexedMenu(title string, items []string) (int, string) {
	if title != "" {
		fmt.Println(title)
	}

	for menuIndex, each := range items {
		if len(items) > 9 {
			fmt.Printf("%02d - %s\n", menuIndex+1, each)
		} else {
			fmt.Printf("%d - %s\n", menuIndex+1, each)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	menuIndex := 0
	for {
		fmt.Print("Enter menu index: ")
		menuIndexStr, _ := reader.ReadString('\n')
		menuIndexStr = strings.TrimSuffix(menuIndexStr, "\n")
		menuIndex, _ = strconv.Atoi(menuIndexStr)

		if menuIndex >= 1 && menuIndex <= len(items) {
			return menuIndex, items[menuIndex-1]
		}
		fmt.Printf("Please enter menu index in the range [%d-%d]\n", 1, len(items))
	}
}
