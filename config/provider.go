package config

import (
	"github.com/spf13/viper"
)

// App struct of app
type App struct {
	Name        string
	Description string // description of app
	Root        string // root path of app
	Prefix      string // prefix key of key
	Tmpl        string // .env.tmpl
	Dist        string // .env
	Cmd         string // Cmd "nginx -s reload"
}

func valueOfMap(key string, m map[string]interface{}, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

// GetApps get app config from .econfig
func GetApps() []*App {
	apps := viper.GetStringMap("apps")
	var appList []*App
	for appName, appProps := range apps {
		appPropsMap := appProps.(map[string]interface{})
		app := &App{
			Name:        appName,
			Description: valueOfMap("description", appPropsMap, ""),
			Root:        valueOfMap("root", appPropsMap, ""),
			Prefix:      valueOfMap("prefix", appPropsMap, ""),
			Tmpl:        valueOfMap("tmpl", appPropsMap, ".env.tmpl"),
			Dist:        valueOfMap("dist", appPropsMap, ".env"),
			Cmd:         valueOfMap("cmd", appPropsMap, ""),
		}
		appList = append(appList, app)
	}
	return appList
}
