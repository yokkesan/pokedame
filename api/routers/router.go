package routers

import (
	"api-generated/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/health", &controllers.HealthController{}, "get:Get")
}
