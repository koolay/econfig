package admin

import (
	"fmt"

	"github.com/iris-contrib/middleware/cors"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/websocket"
)

type WebServer struct {
	setting Setting
}

type Setting struct {
	Port      int
	Username  string
	Password  string
	SecretKey string
}

func New(config Setting) *WebServer {
	return &WebServer{
		setting: config,
	}

}

func (w *WebServer) Start() {

	app := iris.New()
	app.Adapt(httprouter.New())
	app.Adapt(iris.DevLogger())
	// cross domain
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowedHeaders:   []string{"DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	app.Use(crs)

	//jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
	//ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	//return []byte(w.setting.SecretKey), nil
	//},
	//SigningMethod: jwt.SigningMethodHS256,
	//})
	//app.Use(jwtHandler)
	ws := websocket.New(websocket.Config{
		// the path which the websocket client should listen/registered to,
		Endpoint: "/log",
		// the client-side javascript static file path
		// which will be served by Iris.
		// default is /iris-ws.js
		// if you change that you have to change the bottom of templates/client.html
		// script tag:
		ClientSourcePath: "/iris-ws.js",
		//
		// Set the timeouts, 0 means no timeout
		// websocket has more configuration, go to ../../config.go for more:
		// WriteTimeout: 0,
		// ReadTimeout:  0,
		// by-default all origins are accepted, you can change this behavior by setting:
		// CheckOrigin: (r *http.Request ) bool {},
		//
		//
		// IDGenerator used to create (and later on, set)
		// an ID for each incoming websocket connections (clients).
		// The request is an argument which you can use to generate the ID (from headers for example).
		// If empty then the ID is generated by DefaultIDGenerator: randomString(64):
		// IDGenerator func(ctx *iris.Context) string {},
	})
	app.Adapt(ws)
	view := &View{WebServer: w}
	registerRoutes(app, view)
	ws.OnConnection(view.WebSocketHandle)

	app.Listen(fmt.Sprintf(":%d", w.setting.Port))
	//iris.ListenTLSAuto(fmt.Sprintf(":%d", port))
}
