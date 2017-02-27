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

package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/koolay/econfig/cmd"
	"github.com/koolay/econfig/config"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

func main() {
	initializeConfig()
	cmd.Execute()
}

func initLogging(verbose bool) {
	jww.SetLogFile("econfig.log")
	if verbose {
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelInfo)
	} else {
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelError)
	}
}

func initializeConfig() {
	initLogging(config.GlobalFlags.Verbose)
	viper.SetConfigType("toml")
	if config.GlobalFlags.CfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(config.GlobalFlags.CfgFile)
	}

	jww.INFO.Printf("verbose: %v", config.GlobalFlags.Verbose)
	viper.SetConfigName(".econfig") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read in environment variables that match
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
