// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/koolay/econfig/app"
	"github.com/koolay/econfig/config"
	"github.com/koolay/econfig/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveFlag *config.ServeFlag

// ServeCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as a serve",
	Long:  `Run as a serve`,
	Run: func(cmd *cobra.Command, args []string) {
		context.Logger = config.NewLogger(context.Flags.Global)
		context.Logger.INFO.Println(viper.Get("apps.myapp"))
		context.Logger.INFO.Println("serve called")
		cfg := &app.GeneratorConfig{Interval: serveFlag.Interval}
		serveCfg := &app.ServeConfig{}
		serveCfg.Bind = serveFlag.Bind
		serveCfg.Advertise = serveFlag.Advertise
		serveCfg.RPCAddr = serveFlag.RPCAddr
		serveCfg.RPCAuth = serveFlag.RPCAuth
		serveCfg.HttpPort = serveFlag.HttpPort
		serveCfg.Interval = serveFlag.Interval
		serveCfg.Join = serveFlag.Join
		serveCfg.Node = serveFlag.Node
		if serveCfg.Node == "" {
			hostname, err := os.Hostname()
			if err != nil {
				context.Logger.FATAL.Panicf("Error determining hostname: %s", err)
			}
			serveCfg.Node = hostname
		}

		if gen, err := app.NewGenerator(cfg); err == nil {
			go func() {
				gen.SyncLoop()
			}()

		} else {
			context.Logger.FATAL.Panic(err)
		}
		c := app.NewSerfClient(serveCfg)
		if err := c.StartCluster(); err != nil {
			context.Logger.FATAL.Fatal(err)
		}

		// rpcClient, err := app.NewRPCClient(serveCfg.RPCAddr, serveCfg.RPCAuth)
		// if err != nil {
		// context.Logger.FATAL.Fatalf("Create RPC failed. %v", err)
		// }
		//
		// go func() {
		// for {
		// context.Logger.INFO.Panicln("rpc call")
		// err = rpcClient.UserEvent("debug", []byte("hi, rpc"), false)
		// if err != nil {
		// context.Logger.ERROR.Printf("rpc send failed. %v", err)
		// }
		// time.Sleep(5 * time.Second)
		// }
		// }()

		signalChan := make(chan os.Signal, 1)
		doneChan := make(chan bool)
		errChan := make(chan error, 10)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case err := <-errChan:
				context.Logger.ERROR.Printf("%v", err)
			case s := <-signalChan:
				context.Logger.INFO.Printf("Captured %v. Exiting...", s)
				close(doneChan)
			case <-doneChan:
				os.Exit(0)
			}
		}

	},
}

func init() {
	EConfigCmd.AddCommand(ServeCmd)
	serveFlag = config.NewServeFlag(ServeCmd.PersistentFlags())
}
