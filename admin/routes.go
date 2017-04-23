package admin

import (
	"gopkg.in/kataras/iris.v6"
)

func registerAPIRoutes(app *iris.Router, view *View) {
	app.Get("/apps", view.GetApps)
	app.Get("/app/:app", view.GetApp)
	app.Post("/app/:app/sync", view.ExecuteSync)
	app.Get("/app/:app/tmp-items", view.GetAllTmpItems)
	app.Get("/app/:app/store-items", view.GetAllStoredItems)
	app.Get("/app/:app/items", view.AllItems)
	app.Get("/app/:app/item/:key", view.GetItem)
	app.Put("/app/:app/item", view.SetItem)
}
