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
	"fmt"
	"os"

	"github.com/koolay/econfig/config"
	"github.com/koolay/econfig/context"
	"github.com/spf13/cobra"
)

var (
	gf *config.GlobalFlag
)

// EConfigCmd represents the base command when called without any subcommands
var EConfigCmd = &cobra.Command{
	Use:   "econfig",
	Short: "generate dotenv file",
	Long:  `A tool for generating dotenv file from data source.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root run")
	},
}

func Execute() {
	if err := EConfigCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(preInit)
	initRootPersistenFlags()
}

// initConfig reads in config file and ENV variables if set.
func preInit() {
	initConfig()
}

func initRootPersistenFlags() {
	gf = &config.GlobalFlag{}
	EConfigCmd.PersistentFlags().StringVar(&gf.LogFile, "logfile", "", "output log to the logfile. By default output to stdout")
	EConfigCmd.PersistentFlags().StringVarP(&gf.CfgFile, "config", "c", "", "config file")
	EConfigCmd.PersistentFlags().StringVar(&gf.Backend, "backend", "redis", "Backend storage to use, redis, mysql, postgres, etcd, consul.")
	EConfigCmd.PersistentFlags().StringVar(&gf.App, "app", "", "command work only this app")
	EConfigCmd.PersistentFlags().BoolVarP(&gf.Verbose, "verbose", "v", false, "if ouput detail log")
}
func initConfig() {
	context.Flags.Global = gf
	config.LoadConfig(gf)
	context.Logger = config.NewLogger(context.Flags.Global)
}
