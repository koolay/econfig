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
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/koolay/econfig/cmd"
	"github.com/koolay/econfig/config"
	"github.com/koolay/econfig/context"
	"github.com/spf13/pflag"
)

// WordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

func main() {
	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	context.Flags.Global = config.NewGlobalFlag(pflag.CommandLine)
	context.Flags.Serve = config.NewServeFlag(pflag.CommandLine)
	pflag.Parse()
	config.LoadConfig(context.Flags.Global)
	context.Logger = config.NewLogger(context.Flags.Global)

	if err := cmd.EConfigCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
