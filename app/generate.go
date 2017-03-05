// Package app provides ...
package app

import (
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/koolay/econfig/config"
	"github.com/koolay/econfig/context"
	"github.com/koolay/econfig/dotfile"
	"github.com/koolay/econfig/store"
)

type Generator struct {
	mu sync.Mutex
	wg sync.WaitGroup
}

func NewGenerator() (*Generator, error) {
	gen := &Generator{}
	return gen, nil
}

func (gen *Generator) processKey(prefix string, key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return strings.ToLower(fmt.Sprintf("%s_%s", prefix, key))
	} else {
		return strings.ToLower(key)
	}
}

func (gen *Generator) Sync() {
	gen.mu.Lock()
	defer gen.mu.Unlock()
	appList := config.GetApps()
	for _, app := range appList {
		gen.wg.Add(1)
		go func(app *config.App) {
			defer gen.wg.Done()
			context.Logger.INFO.Println("start sync app:", app.Name)
			tmplFilePath := path.Join(app.Root, app.Tmpl)
			// destFilePath := path.Join(app.Root, app.Dist)

			// load apps
			// for: template file , loop keys, if exists key from store,
			if itemMap, err := dotfile.ReadEnvFile(tmplFilePath); err == nil {
				keys := []string{}
				for k, _ := range itemMap {
					keys = append(keys, gen.processKey(app.Prefix, k))
					// context.Logger.DEBUG.Printf("key: %v, v: %v, des: %v \n", k, item.Value, item.Comment)
				}
				storage := store.NewStorage()
				if items, err := storage.GetItems(keys); err == nil {
					context.Logger.INFO.Println(items)
				} else {
					context.Logger.ERROR.Printf("fetch data failed. %v", err)
				}
			} else {
				context.Logger.ERROR.Printf("errmsg: %s. parse file: %s failed.", err.Error(), tmplFilePath)
			}
			context.Logger.INFO.Println("sync app:", app.Name, "completed")
		}(app)

	}
	gen.wg.Wait()
	context.Logger.INFO.Println("complete")
}

func (gen *Generator) SyncLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		<-t.C
		context.Logger.INFO.Println("syncLoop called")
		gen.Sync()
	}
}
