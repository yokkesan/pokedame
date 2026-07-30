package controllers

import (
	"context"
	"errors"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"

	"api-generated/database"
	"api-generated/models"
)

type AdminPokemonController struct {
	beego.Controller
}

type AdminPokemonIndexViewData struct {
	PokemonSpecies []models.PokemonSpecies
	ErrorMessage   string
}

type AdminPokemonDetailViewData struct {
	PokemonSpecies *models.PokemonSpecies
	PokemonForms   []models.PokemonForm
	PokemonAssets  []models.PokemonAsset
	ErrorMessage   string
}

func (c *AdminPokemonController) Index() {
	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		5*time.Second,
	)
	defer cancel()

	db, err := database.Get()
	if err != nil {
		c.Data["ViewData"] = AdminPokemonIndexViewData{
			PokemonSpecies: make([]models.PokemonSpecies, 0),
			ErrorMessage:   "データベースへ接続できませんでした。",
		}
		c.TplName = "admin/pokemon.tpl"
		return
	}

	species, err := models.FindAllPokemonSpecies(
		ctx,
		db,
	)

	viewData := AdminPokemonIndexViewData{
		PokemonSpecies: species,
	}

	if err != nil {
		viewData.PokemonSpecies =
			make([]models.PokemonSpecies, 0)

		viewData.ErrorMessage =
			"ポケモン種族の一覧を取得できませんでした。"
	}

	c.Data["ViewData"] = viewData
	c.TplName = "admin/pokemon.tpl"
}

func (c *AdminPokemonController) Detail() {
	speciesID, err := strconv.ParseInt(
		c.Ctx.Input.Param(":id"),
		10,
		64,
	)

	if err != nil || speciesID <= 0 {
		c.Abort("404")
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		5*time.Second,
	)
	defer cancel()

	species, err := models.FindPokemonSpeciesByID(
		ctx,
		speciesID,
	)

	if errors.Is(
		err,
		models.ErrPokemonSpeciesNotFound,
	) {
		c.Abort("404")
		return
	}

	if err != nil {
		c.Data["ViewData"] = AdminPokemonDetailViewData{
			PokemonForms:  make([]models.PokemonForm, 0),
			PokemonAssets: make([]models.PokemonAsset, 0),
			ErrorMessage:  "ポケモン種族の詳細を取得できませんでした。",
		}
		c.TplName = "admin/pokemon_detail.tpl"
		return
	}

	db, err := database.Get()
	if err != nil {
		c.Data["ViewData"] = AdminPokemonDetailViewData{
			PokemonSpecies: species,
			PokemonForms:   make([]models.PokemonForm, 0),
			PokemonAssets:  make([]models.PokemonAsset, 0),
			ErrorMessage:   "データベースへ接続できませんでした。",
		}
		c.TplName = "admin/pokemon_detail.tpl"
		return
	}

	forms, err := models.FindPokemonFormsBySpeciesID(
		ctx,
		db,
		speciesID,
	)

	viewData := AdminPokemonDetailViewData{
		PokemonSpecies: species,
		PokemonForms:   forms,
		PokemonAssets:  make([]models.PokemonAsset, 0),
	}

	if err != nil {
		viewData.PokemonForms =
			make([]models.PokemonForm, 0)

		viewData.ErrorMessage =
			"フォーム一覧を取得できませんでした。"

		c.Data["ViewData"] = viewData
		c.TplName = "admin/pokemon_detail.tpl"
		return
	}

	for _, form := range forms {
		assets, err := models.FindPokemonAssetsByFormID(
			ctx,
			db,
			form.ID,
		)

		if err != nil {
			viewData.ErrorMessage =
				"素材一覧を取得できませんでした。"
			break
		}

		viewData.PokemonAssets = append(
			viewData.PokemonAssets,
			assets...,
		)
	}

	c.Data["ViewData"] = viewData
	c.TplName = "admin/pokemon_detail.tpl"
}
