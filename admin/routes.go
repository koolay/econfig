package admin

import (
	"gopkg.in/kataras/iris.v6"
)

func registerRoutes(app *iris.Framework, view *View) {

	app.Get("/", view.ServeStatic)
	app.Get("/static/*file", view.ServeStatic)
	app.Get("/view/*file", view.ServeStatic)

	//login
	app.Get("/home", view.Home)
	app.Post("/api/login", view.Login)
	app.Post("/api/exec", view.Execute)
	app.Get("/api/apps", view.GetApps)
	app.Get("/api/app/:app", view.GetApp)
	app.Get("/api/app/:app/tmp-items", view.GetAllTmpItems)
	app.Get("/api/app/:app/store-items", view.GetAllStoredItems)
	app.Get("/api/app/:app/items", view.AllItems)
	app.Get("/api/app/:app/item/:key", view.GetItem)
}
