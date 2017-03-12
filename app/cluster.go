// Package app
package app

import (
	"fmt"
	"net"
	"sync"

	"github.com/pkg/errors"

	"github.com/hashicorp/serf/serf"
	"github.com/koolay/econfig/config"
	"github.com/koolay/econfig/context"
)

const (
	DefaultBindPort = 7946
)

type ServeConfig struct {
	config.ServeFlag
}

// Client serf cluster client
type Client struct {
	// eventCh is used to receive events from the
	// serf cluster in the datacenter
	cfg     *ServeConfig
	eventCh chan serf.Event

	// eventHandlers is the registered handlers for events
	eventHandlers     map[EventHandler]struct{}
	eventHandlerList  []EventHandler
	eventHandlersLock sync.Mutex

	serf *serf.Serf
	// shutdownCh is used for shutdowns
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// BindAddrParts returns the parts of the BindAddr that should be
// used to configure Serf.
func (c *ServeConfig) AddrParts(address string) (string, int, error) {
	checkAddr := address

START:
	_, _, err := net.SplitHostPort(checkAddr)
	if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
		checkAddr = fmt.Sprintf("%s:%d", checkAddr, DefaultBindPort)
		goto START
	}
	if err != nil {
		return "", 0, err
	}

	// Get the address
	addr, err := net.ResolveTCPAddr("tcp", checkAddr)
	if err != nil {
		return "", 0, err
	}

	return addr.IP.String(), addr.Port, nil
}

func NewSerfClient(cfg *ServeConfig) *Client {
	return &Client{cfg: cfg}
}

// setupSerf is used to setup and initialize a Serf
func (c *Client) initSerfConfig() (*serf.Config, error) {

	context.Logger.DEBUG.Println("setup serf cluster")
	conf := serf.DefaultConfig()
	conf.NodeName = c.cfg.Node

	var advertiseIP string
	var advertisePort int = DefaultBindPort
	var err error

	bindIP, bindPort, err := c.cfg.AddrParts(c.cfg.Bind)
	if err != nil {
		context.Logger.FATAL.Panicf("Invalid bind address: %s", err)
	}

	if c.cfg.Advertise != "" {
		advertiseIP, advertisePort, err = c.cfg.AddrParts(c.cfg.Advertise)
		if err != nil {
			return nil, errors.Errorf("Invalid advertise address: %s", err)
		}
	}

	conf.MemberlistConfig.AdvertiseAddr = advertiseIP
	conf.MemberlistConfig.AdvertisePort = advertisePort

	conf.MemberlistConfig.BindAddr = bindIP
	conf.MemberlistConfig.BindPort = bindPort
	return conf, nil
}

// Start is used to initiate the event listeners. It is separate from
// create so that there isn't a race condition between creating the
// agent and registering handlers
func (c *Client) StartCluster() error {

	serfConfig, err := c.initSerfConfig()
	if err != nil {
		return err
	}

	context.Logger.INFO.Printf("agent starting")
	// Create serf first
	cluster, err := serf.Create(serfConfig)
	if err != nil {
		return errors.WithMessage(err, "Couldn't create cluster")
	}

	context.Logger.INFO.Printf("Setup cluster, bind: %s:%d \n", serfConfig.MemberlistConfig.BindAddr, serfConfig.MemberlistConfig.BindPort)

	if serfConfig.MemberlistConfig.AdvertiseAddr != "" {
		context.Logger.INFO.Printf("advertise: %s:%d \n", serfConfig.MemberlistConfig.AdvertiseAddr, serfConfig.MemberlistConfig.AdvertisePort)
	}

	context.Logger.INFO.Printf("Node name: %s \n", serfConfig.NodeName)
	context.Logger.INFO.Printf("Http port: %d \n", c.cfg.HttpPort)

	c.serf = cluster
	// start event loop
	go c.eventLoop()

	// join cluster
	if c.cfg.Join != "" {
		context.Logger.DEBUG.Println(" Joinning cluster: ", c.cfg.Join)
		n, err := cluster.Join([]string{c.cfg.Join}, true)
		if err != nil {
			context.Logger.FATAL.Fatalf("Couldn't join cluster, starting own: %v\n", err)
		}
		if n > 0 {
			context.Logger.INFO.Printf("agent: joined: %d nodes", n)
		}
		if err != nil {
			context.Logger.WARN.Printf("agent: error joining: %v", err)
		}
	}
	return nil

}

// eventLoop listens to events from Serf and fans out to event handlers
func (c *Client) eventLoop() {
	serfShutdownCh := c.serf.ShutdownCh()
	for {
		select {
		case e := <-c.eventCh:
			context.Logger.INFO.Printf("agent: Received event: %s", e.String())
			c.eventHandlersLock.Lock()
			handlers := c.eventHandlerList
			c.eventHandlersLock.Unlock()
			for _, eh := range handlers {
				eh.HandleEvent(e)
			}

		case <-serfShutdownCh:
			context.Logger.WARN.Printf("agent: Serf shutdown detected, quitting")
			c.Shutdown()
			return

		case <-c.shutdownCh:
			return
		}
	}
}

// Shutdown closes this agent and all of its processes. Should be preceded
// by a Leave for a graceful shutdown.
func (c *Client) Shutdown() error {
	c.shutdownLock.Lock()
	defer c.shutdownLock.Unlock()

	if c.shutdown {
		return nil
	}

	if c.serf == nil {
		goto EXIT
	}

	context.Logger.INFO.Println("agent: requesting serf shutdown")
	if err := c.serf.Shutdown(); err != nil {
		return err
	}

EXIT:
	context.Logger.INFO.Println("agent: shutdown complete")
	c.shutdown = true
	close(c.shutdownCh)
	return nil
}
