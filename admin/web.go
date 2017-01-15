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
	Port      int    `web port`
	Username  string `admin username`
	Password  string `admin password`
	SecretKey string `jwt secretKey`
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
	app.Get("/api/projects", jwtMDW.Serve, view.GetProjects)
	app.Get("/api/project/:projectName", jwtMDW.Serve, view.GetProject)
	app.Get("/api/project/:projectName/item/:key", jwtMDW.Serve, view.GetItem)
	app.Delete("/api/project/:projectName/item/:key", jwtMDW.Serve, view.DeleteItem)
	app.Get("/api/project/:projectName/items", jwtMDW.Serve, view.GetItems)
	app.Post("/api/project/:projectName/items", jwtMDW.Serve, view.SetItem)
	//tmpl
	app.Get("/api/project/:projectName/tmpl/:filepath", jwtMDW.Serve, view.GetTemplates)
	app.Websocket.OnConnection(view.WebSocketHandle)

	app.Listen(fmt.Sprintf(":%d", w.setting.Port))
	//iris.ListenTLSAuto(fmt.Sprintf(":%d", port))
}
