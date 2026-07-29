package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"api-generated/database"
	"api-generated/models"

	beego "github.com/beego/beego/v2/server/web"
)

const (
	pokemonSpeciesRequestTimeout = 5 * time.Second
	maxJSONRequestBodySize       = 1 << 20
)

type PokemonSpeciesController struct {
	beego.Controller
}

type pokemonSpeciesErrorResponse struct {
	Message string `json:"message"`
}

func (c *PokemonSpeciesController) Create() {
	var request models.CreatePokemonSpeciesRequest

	if err := decodePokemonSpeciesJSON(c, &request); err != nil {
		c.writePokemonSpeciesError(
			http.StatusBadRequest,
			"リクエスト形式が不正です。",
		)
		return
	}

	db, err := database.Get()
	if err != nil {
		c.writePokemonSpeciesError(
			http.StatusInternalServerError,
			"データベースへ接続できませんでした。",
		)
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		pokemonSpeciesRequestTimeout,
	)
	defer cancel()

	pokemon, err := models.CreatePokemonSpecies(ctx, db, request)
	if err != nil {
		c.handlePokemonSpeciesCreateError(err)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = pokemon
	c.ServeJSON()
}

func (c *PokemonSpeciesController) List() {
	db, err := database.Get()
	if err != nil {
		c.writePokemonSpeciesError(
			http.StatusInternalServerError,
			"データベースへ接続できませんでした。",
		)
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		pokemonSpeciesRequestTimeout,
	)
	defer cancel()

	pokemonList, err := models.FindAllPokemonSpecies(ctx, db)
	if err != nil {
		c.writePokemonSpeciesError(
			http.StatusInternalServerError,
			"ポケモン一覧を取得できませんでした。",
		)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = pokemonList
	c.ServeJSON()
}

func (c *PokemonSpeciesController) handlePokemonSpeciesCreateError(err error) {
	switch {
	case errors.Is(err, models.ErrInvalidNationalDexNumber):
		c.writePokemonSpeciesError(
			http.StatusBadRequest,
			"全国図鑑番号が不正です。",
		)
	case errors.Is(err, models.ErrInvalidPokemonNameJA):
		c.writePokemonSpeciesError(
			http.StatusBadRequest,
			"日本語名が不正です。",
		)
	case errors.Is(err, models.ErrInvalidPokemonNameEN):
		c.writePokemonSpeciesError(
			http.StatusBadRequest,
			"英語名が不正です。",
		)
	case errors.Is(err, models.ErrInvalidPokemonSlug):
		c.writePokemonSpeciesError(
			http.StatusBadRequest,
			"スラッグが不正です。",
		)
	case errors.Is(err, models.ErrDuplicatePokemonSpecies):
		c.writePokemonSpeciesError(
			http.StatusConflict,
			"同じ図鑑番号またはスラッグが既に登録されています。",
		)
	default:
		c.writePokemonSpeciesError(
			http.StatusInternalServerError,
			"ポケモンを登録できませんでした。",
		)
	}
}

func (c *PokemonSpeciesController) writePokemonSpeciesError(
	status int,
	message string,
) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = pokemonSpeciesErrorResponse{
		Message: message,
	}
	c.ServeJSON()
}

func decodePokemonSpeciesJSON(
	c *PokemonSpeciesController,
	destination any,
) error {
	reader := io.LimitReader(
		c.Ctx.Request.Body,
		maxJSONRequestBodySize,
	)

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("multiple JSON values are not allowed")
	}

	return nil
}
