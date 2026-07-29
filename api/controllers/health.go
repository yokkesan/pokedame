package controllers

import (
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type HealthController struct {
	beego.Controller
}

type HealthResponse struct {
	Status string `json:"status"`
}

func (c *HealthController) Get() {
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = HealthResponse{
		Status: "ok",
	}
	c.ServeJSON()
}
