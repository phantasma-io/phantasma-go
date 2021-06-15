package rpc_test

import (
	"testing"

	"github.com/phantasma.io/phantasma-go/pkg/rpc"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := rpc.NewRPCMainnet()
	assert.NotNil(t, client)
}

func TestAccount(t *testing.T) {
	client := rpc.NewRPCMainnet()
	account, err := client.GetAccount("P2KA7yzB3uUncuAqP6tLut27iTKAC6ZTnAVM4myUuG57oQP")
	assert.Nil(t, err)
	assert.NotNil(t, account)
}

//func TestAccounts(t *testing.T) {
//	client := rpc.NewRPCMainnet()
//	accounts, err := client.GetAccounts("P2KA7yzB3uUncuAqP6tLut27iTKAC6ZTnAVM4myUuG57oQP,P2K4M8KVTqg1eKTuvtp5hETGCNkzjhaRtbJJQN97qJVpAZz")
//	assert.Nil(t, err)
//	fmt.Println(accounts)
//	assert.True(t, len(accounts) == 2)
//	assert.NotNil(t, accounts)
//}
