package config

import (
	"time"

	"github.com/spf13/pflag"
)

// GlobalFlag all flags of command
type GlobalFlag struct {
	// full path of config file
	CfgFile string

	// redis, mysql, etcd, consul etc..
	Store string

	// process which app
	App string

	Verbose bool
}

type ServeFlag struct {
	Interval time.Duration
}

func NewServeFlag(fs *pflag.FlagSet) *ServeFlag {
	gf := &ServeFlag{}
	gf.Interval = 60 * time.Second
	fs.DurationVar(&gf.Interval, "interval", gf.Interval, "seconds for polling interval(default 60)")
	return gf
}

func NewGlobalFlag(fs *pflag.FlagSet) *GlobalFlag {
	gf := &GlobalFlag{}
	fs.StringVarP(&gf.CfgFile, "config", "c", "", "config file")
	fs.StringVar(&gf.Store, "store", "redis", "configuration where to store, redis,mysql,etcd,consul")
	fs.StringVar(&gf.App, "app", "", "command work only this app")
	fs.BoolVarP(&gf.Verbose, "verbose", "v", false, "if ouput detail log")
	return gf
}
