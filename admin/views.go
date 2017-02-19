package admin

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

type View struct {
	WebServer *WebServer
}

func (v *View) Execute(ctx *iris.Context) {
}

func (v *View) WebSocketHandle(c iris.WebsocketConnection) {
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
	path := ctx.PathString()
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
			ctx.JSON(iris.StatusOK, iris.Map{"result": true, "token": tokenString})
		} else {
			ctx.JSON(iris.StatusOK, iris.Map{"result": false, "msg": err.Error()})
		}
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"result": false, "msg": "username or password incorrect"})
	}
}
