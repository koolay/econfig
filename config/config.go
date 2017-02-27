package config

// GlobalFlags hold flags
var GlobalFlags Flags

// Flags all flags of command
type Flags struct {
	// full path of config file
	CfgFile string

	// redis, mysql, etcd, consul etc..
	Store string

	// process which app
	App string

	Verbose bool
}

func init() {
}
