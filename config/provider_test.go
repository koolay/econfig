package config_test

import (
	"testing"

	"github.com/koolay/econfig/config"
	"github.com/stretchr/testify/assert"
)

var configTmp = `
[flags]
store = "redis"

[web]
username = "admin"
password = "123"

[apps]
[apps.myapp]
description = "myapp"
root = "/data/www/myapp"
prefix = ""
tmpl = ".env.example"
dist = ".env.dist"
cmd = "nginx -s reload"

[apps.myapp2]
description = "myapp2"
root = "/data/www/myapp2"
cmd = ""

[store]
[store.redis]
host = "127.0.0.1"
port = 6379
db = 0
password = ""

[store.mysql]
host = "127.0.0.1"
port = 6379
database = "econfig"
username = "root"
password = "dev"
`

func TestGetApps(t *testing.T) {
	apps := config.GetApps()
	assert.Equal(t, len(apps), 2)
	assert.Contains(t, []string{"myapp", "myapp2"}, apps[0].Name, "failed contains name")
	assert.Contains(t, []string{".env.tmpl", ".env.example"}, apps[0].Tmpl, "failed contains tmpl")
	assert.Contains(t, []string{".env", ".env.dist"}, apps[0].Dist, "faild contains dist")
}
