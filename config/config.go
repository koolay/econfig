package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// FlagValues hold flags
var FlagValues Flags

// Flags all flags of command
type Flags struct {
	// full path of config file
	CfgFile string

	// redis, mysql, etcd, consul etc..
	Store string

	// process which app
	App string
}

func init() {
	viper.SetConfigType("toml")
	if FlagValues.CfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(FlagValues.CfgFile)
	}

	viper.SetConfigName(".econfig") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read in environment variables that match
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
