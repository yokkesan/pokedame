package controllers

import (
	"context"
	"net/http"
	"time"

	"api-generated/database"
	"api-generated/models"

	beego "github.com/beego/beego/v2/server/web"
)

const adminPokemonRequestTimeout = 5 * time.Second

type AdminPokemonController struct {
	beego.Controller
}

type adminPokemonViewData struct {
	PokemonSpecies []models.PokemonSpecies
	ErrorMessage   string
}

func (c *AdminPokemonController) Index() {
	db, err := database.Get()
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["ViewData"] = adminPokemonViewData{
			PokemonSpecies: []models.PokemonSpecies{},
			ErrorMessage:   "データベースへ接続できませんでした。",
		}
		c.TplName = "admin/pokemon.tpl"
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		adminPokemonRequestTimeout,
	)
	defer cancel()

	pokemonSpecies, err := models.FindAllPokemonSpecies(ctx, db)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["ViewData"] = adminPokemonViewData{
			PokemonSpecies: []models.PokemonSpecies{},
			ErrorMessage:   "ポケモン一覧を取得できませんでした。",
		}
		c.TplName = "admin/pokemon.tpl"
		return
	}

	c.Data["ViewData"] = adminPokemonViewData{
		PokemonSpecies: pokemonSpecies,
	}
	c.TplName = "admin/pokemon.tpl"
}
