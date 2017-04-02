package admin

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/koolay/econfig/config"
	"github.com/koolay/econfig/context"
	"github.com/koolay/econfig/dotfile"
	"github.com/koolay/econfig/store"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/websocket"
)

type View struct {
	WebServer *WebServer
}

func (v *View) sendJson(ctx *iris.Context, data interface{}) {
	if data == nil {
		ctx.JSON(200, iris.Map{"code": 0, "msg": "ok"})
	} else {
		ctx.JSON(200, iris.Map{"code": 0, "data": data})
	}
}

func (v *View) sendJsonForError(ctx *iris.Context, code int, err interface{}) {
	switch err.(type) {
	case string:
		ctx.JSON(200, iris.Map{"code": code, "msg": err.(string)})
	case error:
		ctx.JSON(200, iris.Map{"code": code, "msg": err.(error).Error()})
	default:
		if data, ex := json.Marshal(err); ex == nil {
			ctx.JSON(200, iris.Map{"code": code, "msg": data})
		} else {
			ctx.JSON(200, iris.Map{"code": code, "msg": ex.Error()})
		}
	}
}

func (v *View) newStore() (store.Storage, error) {
	return store.NewStorage(context.Flags.Global.Backend)
}

func (v *View) Execute(ctx *iris.Context) {
}

func (v *View) WebSocketHandle(c websocket.Connection) {
	log.Debug("client connet now! ID: %s", c.ID())
	c.Join("confd")
	c.On("log", func(message string) {
		// to all except this connection ->
		//c.To(iris.Broadcast).Emit("chat", "Message from: "+c.ID()+"-> "+message)

		// to the client ->
		//c.Emit("chat", "Message from myself: "+message)

		c.To("confd").Emit("log", "replay from server message!")
		// send the message to the whole room,
		// all connections which are inside this room will receive this message
		//c.To("confd").Emit("chat", "From: "+c.ID()+": "+message)
	})

	c.OnDisconnect(func() {
		log.Debug("Connection with ID: %s has been disconnected!", c.ID())
	})
}

func (v *View) ServeStatic(ctx *iris.Context) {
	path := ctx.Path()
	log.Debug("service path:" + path)

	if path == "/" || (!strings.Contains(path, ".js") && !strings.Contains(path, ".css") && !strings.Contains(path, ".png") && !strings.Contains(path, ".icon") && !strings.Contains(path, ".gif") && !strings.Contains(path, ".ttf") && !strings.Contains(path, ".woff")) {
		path = "index.html"
	}

	path = filepath.Join("web/dist/", path)
	path = strings.Replace(path, "/", string(os.PathSeparator), -1)
	path = strings.TrimPrefix(path, "/")
	if uri, err := url.Parse(path); err == nil {
		path = uri.Path
	} else {
		ctx.Text(iris.StatusInternalServerError, err.Error())
		return
	}

	log.Debug("static path:" + path)
	data, err := Asset(path)
	if err != nil {
		log.Error(err.Error())
		ctx.NotFound()
		return
	}

	ctx.ServeContent(bytes.NewReader(data), path, time.Now(), true)
}

func (v *View) GetItem(ctx *iris.Context) {
	appName := ctx.Param("app")
	key := ctx.Param("key")
	app := config.GetApp(appName)
	if app == nil {
		v.sendJsonForError(ctx, 404, "app not exist")
	}

	if item, err := dotfile.GetConfigItemFromEnvFile(app.GetDestPath(), key); err == nil {
		if item == nil {
			v.sendJsonForError(ctx, 404, "key not exist")
			return
		} else {
			v.sendJson(ctx, item)
			return
		}
	} else {
		v.sendJsonForError(ctx, 500, err.Error())
		return
	}
}

func (v *View) SetItem(ctx *iris.Context) {
	appName := ctx.Param("app")
	key := ctx.PostValue("key")
	value := ctx.PostValue("value")
	if key == "" {
		v.sendJsonForError(ctx, 403, "key not allowed empty")
		return
	}
	app := config.GetApp(appName)
	if app == nil {
		v.sendJsonForError(ctx, 404, "app not exist")
		return
	}
	storeKey := app.GenerateStoreKey(key)
	storage, err := v.newStore()
	if err == nil {
		err = storage.SetItem(storeKey, value)
	}
	if err == nil {
		v.sendJson(ctx, nil)
	} else {
		v.sendJsonForError(ctx, 500, err.Error())
	}
}

func (v *View) AllItems(ctx *iris.Context) {
	appName := ctx.Param("app")
	app := config.GetApp(appName)
	if app == nil {
		v.sendJsonForError(ctx, 404, "app not exist")
		return
	}

	items, err := dotfile.ReadEnvFile(app.GetDestPath())
	if err != nil {
		v.sendJsonForError(ctx, 500, err.Error())
		return
	} else {
		v.sendJson(ctx, items)
		return
	}
}

func (v *View) GetAllTmpItems(ctx *iris.Context) {
	appName := ctx.Param("app")
	app := config.GetApp(appName)
	if app == nil {
		v.sendJsonForError(ctx, 404, "app not exist")
		return
	}
	itemsMap, err := dotfile.ReadEnvFile(app.GetTmplPath())
	if err == nil {
		v.sendJson(ctx, itemsMap)
	} else {
		v.sendJsonForError(ctx, 500, err)
	}
}

func (v *View) GetAllStoredItems(ctx *iris.Context) {
	appName := ctx.Param("app")
	app := config.GetApp(appName)
	if app == nil {
		v.sendJsonForError(ctx, 404, "app not exist")
		return
	}
	storage, err := v.newStore()
	if err == nil {
		itemsMap, err := dotfile.ReadEnvFile(app.GetTmplPath())
		var keys []string
		for key, _ := range itemsMap {
			keys = append(keys, key)
		}
		items, err := storage.GetItems(keys)
		if err == nil {
			v.sendJson(ctx, items)
			return
		}
	} else {
		v.sendJsonForError(ctx, 500, err)
	}
}

func (v *View) GetApps(ctx *iris.Context) {
	apps := config.GetApps()
	v.sendJson(ctx, apps)
}

func (v *View) GetApp(ctx *iris.Context) {
	appName := ctx.Param("app")
	app := config.GetApp(appName)
	if app == nil {
		v.sendJsonForError(ctx, 404, "app not exist")
	} else {
		v.sendJson(ctx, app)
	}
}

func (v *View) Home(ctx *iris.Context) {
	ctx.WriteString("hello")
}

type User struct {
	Username string
	Password string
}

func (v *View) Login(ctx *iris.Context) {

	username := ctx.PostValue("username")
	password := ctx.PostValue("password")
	log.Debug(username + ", pwd:" + password)
	log.Debug("config username:" + v.WebServer.setting.Username)

	if username == v.WebServer.setting.Username && password == v.WebServer.setting.Password {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		if tokenString, err := token.SignedString([]byte(v.WebServer.setting.SecretKey)); err == nil {
			ctx.JSON(iris.StatusOK, iris.Map{"code": 0, "token": tokenString})
		} else {
			ctx.JSON(iris.StatusOK, iris.Map{"code": 500, "msg": err.Error()})
		}
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"result": 403, "msg": "username or password incorrect"})
	}
}
