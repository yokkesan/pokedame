package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"api-generated/database"
	"api-generated/models"

	beego "github.com/beego/beego/v2/server/web"
)

const pokemonFormRequestTimeout = 5 * time.Second

type PokemonFormController struct {
	beego.Controller
}

type pokemonFormErrorResponse struct {
	Message string `json:"message"`
}

func (c *PokemonFormController) Create() {
	pokemonSpeciesID, err := c.getPokemonSpeciesID()
	if err != nil {
		c.writePokemonFormError(
			http.StatusBadRequest,
			"ポケモン種族IDが不正です。",
		)
		return
	}

	var request models.CreatePokemonFormRequest

	if err := decodePokemonFormJSON(c, &request); err != nil {
		c.writePokemonFormError(
			http.StatusBadRequest,
			"リクエスト形式が不正です。",
		)
		return
	}

	db, err := database.Get()
	if err != nil {
		c.writePokemonFormError(
			http.StatusInternalServerError,
			"データベースへ接続できませんでした。",
		)
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		pokemonFormRequestTimeout,
	)
	defer cancel()

	pokemonForm, err := models.CreatePokemonForm(
		ctx,
		db,
		pokemonSpeciesID,
		request,
	)
	if err != nil {
		c.handlePokemonFormCreateError(err)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = pokemonForm
	c.ServeJSON()
}

func (c *PokemonFormController) List() {
	pokemonSpeciesID, err := c.getPokemonSpeciesID()
	if err != nil {
		c.writePokemonFormError(
			http.StatusBadRequest,
			"ポケモン種族IDが不正です。",
		)
		return
	}

	db, err := database.Get()
	if err != nil {
		c.writePokemonFormError(
			http.StatusInternalServerError,
			"データベースへ接続できませんでした。",
		)
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		pokemonFormRequestTimeout,
	)
	defer cancel()

	pokemonForms, err := models.FindPokemonFormsBySpeciesID(
		ctx,
		db,
		pokemonSpeciesID,
	)
	if err != nil {
		c.handlePokemonFormListError(err)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = pokemonForms
	c.ServeJSON()
}

func (c *PokemonFormController) handlePokemonFormCreateError(err error) {
	switch {
	case errors.Is(err, models.ErrInvalidPokemonSpeciesID):
		c.writePokemonFormError(
			http.StatusBadRequest,
			"ポケモン種族IDが不正です。",
		)

	case errors.Is(err, models.ErrPokemonSpeciesNotFound):
		c.writePokemonFormError(
			http.StatusNotFound,
			"指定されたポケモン種族が存在しません。",
		)

	case errors.Is(err, models.ErrInvalidPokemonFormKey):
		c.writePokemonFormError(
			http.StatusBadRequest,
			"フォームキーが不正です。",
		)

	case errors.Is(err, models.ErrInvalidPokemonFormNameJA):
		c.writePokemonFormError(
			http.StatusBadRequest,
			"フォームの日本語名が不正です。",
		)

	case errors.Is(err, models.ErrInvalidPokemonFormNameEN):
		c.writePokemonFormError(
			http.StatusBadRequest,
			"フォームの英語名が不正です。",
		)

	case errors.Is(err, models.ErrDuplicatePokemonForm):
		c.writePokemonFormError(
			http.StatusConflict,
			"同じフォームキーが既に登録されています。",
		)

	case errors.Is(err, models.ErrDuplicateDefaultForm):
		c.writePokemonFormError(
			http.StatusConflict,
			"デフォルトフォームは既に登録されています。",
		)

	default:
		c.writePokemonFormError(
			http.StatusInternalServerError,
			"ポケモンフォームを登録できませんでした。",
		)
	}
}

func (c *PokemonFormController) handlePokemonFormListError(err error) {
	switch {
	case errors.Is(err, models.ErrInvalidPokemonSpeciesID):
		c.writePokemonFormError(
			http.StatusBadRequest,
			"ポケモン種族IDが不正です。",
		)

	case errors.Is(err, models.ErrPokemonSpeciesNotFound):
		c.writePokemonFormError(
			http.StatusNotFound,
			"指定されたポケモン種族が存在しません。",
		)

	default:
		c.writePokemonFormError(
			http.StatusInternalServerError,
			"ポケモンフォーム一覧を取得できませんでした。",
		)
	}
}

func (c *PokemonFormController) getPokemonSpeciesID() (int64, error) {
	value := c.Ctx.Input.Param(":speciesId")

	pokemonSpeciesID, err := strconv.ParseInt(value, 10, 64)
	if err != nil || pokemonSpeciesID <= 0 {
		return 0, models.ErrInvalidPokemonSpeciesID
	}

	return pokemonSpeciesID, nil
}

func (c *PokemonFormController) writePokemonFormError(
	status int,
	message string,
) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = pokemonFormErrorResponse{
		Message: message,
	}
	c.ServeJSON()
}

func decodePokemonFormJSON(
	c *PokemonFormController,
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
