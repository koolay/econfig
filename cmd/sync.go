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

	"github.com/koolay/econfig/app"
	"github.com/koolay/econfig/config"
	"github.com/koolay/econfig/context"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "update config from store immediately",
	Long:  `Update or generate the config file from store immediately.`,
	Run: func(cmd *cobra.Command, args []string) {
		context.Logger = config.NewLogger(context.Flags.Global)
		cfg := &app.GeneratorConfig{}
		if gen, err := app.NewGenerator(cfg); err == nil {
			fmt.Println("sync process")
			gen.Sync(nil)
		} else {
			context.Logger.FATAL.Panic(err)
		}
	},
}

func init() {
	EConfigCmd.AddCommand(syncCmd)
}
