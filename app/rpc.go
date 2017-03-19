package app

import (
	"github.com/hashicorp/serf/client"
)

// RPCClient returns a new Serf RPC client with the given address.
func NewRPCClient(addr, auth string) (*client.RPCClient, error) {
	config := client.Config{Addr: addr, AuthKey: auth}
	return client.ClientFromConfig(&config)
}
