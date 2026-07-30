package routers

import (
	"api-generated/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.SetStaticPath(
		"/pokemon-assets",
		"/app/storage/pokemon-assets",
	)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/health", &controllers.HealthController{}, "get:Get")

	beego.Router(
		"/api/admin/pokemon-species",
		&controllers.PokemonSpeciesController{},
		"get:List;post:Create",
	)

	beego.Router(
		"/api/admin/pokemon-species/:speciesId/forms",
		&controllers.PokemonFormController{},
		"get:List;post:Create",
	)

	beego.Router(
		"/api/admin/pokemon-forms/:formId/assets",
		&controllers.PokemonAssetController{},
		"get:List;post:Create",
	)

	beego.Router(
		"/api/admin/pokemon-assets/:assetId",
		&controllers.PokemonAssetController{},
		"delete:Delete",
	)

	beego.Router(
		"/admin/pokemon",
		&controllers.AdminPokemonController{},
		"get:Index",
	)

	beego.Router(
		"/admin/pokemon/:id",
		&controllers.AdminPokemonController{},
		"get:Detail",
	)
}
