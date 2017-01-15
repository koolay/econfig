package admin

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kelseyhightower/confd/log"
)

type View struct {
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
	go func() {
		lq := log.GetLogQueue()
		for {
			logMessage := lq.GetLatest()
			if logMessage != "" {
				c.To("confd").Emit("log", logMessage)
			}
		}
	}()

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
