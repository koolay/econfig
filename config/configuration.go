package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	supportedBackends map[string]int = map[string]int{"redis": 1, "mysql": 1, "postgres": 1, "pg": 1}
)

// App struct of app
type App struct {
	Name        string
	Description string // description of app
	Root        string // root path of app
	Prefix      string // prefix key of key
	Tmpl        string // .env.tmpl
	Dest        string // .env
	Cmd         string // Cmd "nginx -s reload"
	Owner       string // who owner the dist file
}

func ValueOfMap(key string, m map[string]interface{}, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

func GetBackends(backend string) (map[string]interface{}, error) {

	backend = strings.ToLower(backend)
	if _, ok := supportedBackends[backend]; !ok {
		return nil, errors.New("not supported backend:" + backend)
	}
	backends := viper.GetStringMap("backends")

	for name, props := range backends {
		if backend == strings.ToLower(name) {
			return props.(map[string]interface{}), nil
		}
	}
	return nil, errors.New("not found configuration of the backend:" + backend)
}

// GetApps get app config from .econfig
func GetApps() []*App {
	apps := viper.GetStringMap("apps")
	var appList []*App
	for appName, appProps := range apps {
		appPropsMap := appProps.(map[string]interface{})
		app := &App{
			Name:        appName,
			Description: ValueOfMap("description", appPropsMap, ""),
			Root:        ValueOfMap("root", appPropsMap, ""),
			Prefix:      ValueOfMap("prefix", appPropsMap, ""),
			Tmpl:        ValueOfMap("tmpl", appPropsMap, ".env.tmpl"),
			Dest:        ValueOfMap("dest", appPropsMap, ".env"),
			Cmd:         ValueOfMap("cmd", appPropsMap, ""),
		}
		appList = append(appList, app)
	}
	return appList
}

func LoadConfig(gf *GlobalFlag) {
	viper.SetConfigType("toml")
	if gf.CfgFile != "" {
		viper.SetConfigFile(gf.CfgFile)
	}

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
