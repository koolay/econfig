package admin

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/cors"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kelseyhightower/confd/resource/template"
)

type WebServer struct {
	templateConfig template.Config
	setting        Setting
}

type Setting struct {
	Port      int
	Username  string
	Password  string
	SecretKey string
}

func New(templateConfig template.Config, config Setting) *WebServer {
	return &WebServer{
		templateConfig: templateConfig,
		setting:        config,
	}

}

func (w *WebServer) Start() {

	// cross domain
	crs := cors.New(cors.Options{
		//AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowedHeaders:   []string{"DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})
	crs.Log = iris.Logger

	// jwt middleware
	jwtMDW := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(w.setting.SecretKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	config := iris.Configuration{Charset: "UTF-8", Gzip: true, DisablePathEscape: true}
	app := iris.New(config)
	app.Use(crs)
	app.Config.Websocket.Endpoint = "/log"

	//app.Favicon("./favicon.ico")
	view := &View{WebServer: w}

	//service static file
	app.Get("/", view.ServeStatic)
	app.Get("/static/*file", view.ServeStatic)
	app.Get("/view/*file", view.ServeStatic)

	//login
	app.Post("/api/login", view.Login)
	app.Post("/api/exec", jwtMDW.Serve, view.Execute)
	app.Websocket.OnConnection(view.WebSocketHandle)

	app.Listen(fmt.Sprintf(":%d", w.setting.Port))
	//iris.ListenTLSAuto(fmt.Sprintf(":%d", port))
}
