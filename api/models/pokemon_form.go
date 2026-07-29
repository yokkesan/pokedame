package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	maxPokemonFormKeyLength  = 100
	maxPokemonFormNameLength = 100
)

var (
	ErrInvalidPokemonSpeciesID  = errors.New("invalid pokemon species id")
	ErrPokemonSpeciesNotFound   = errors.New("pokemon species not found")
	ErrInvalidPokemonFormKey    = errors.New("invalid pokemon form key")
	ErrInvalidPokemonFormNameJA = errors.New("invalid pokemon form japanese name")
	ErrInvalidPokemonFormNameEN = errors.New("invalid pokemon form english name")
	ErrDuplicatePokemonForm     = errors.New("duplicate pokemon form")
	ErrDuplicateDefaultForm     = errors.New("default pokemon form already exists")
)

var pokemonFormKeyPattern = regexp.MustCompile(
	`^[a-z0-9]+(?:-[a-z0-9]+)*$`,
)

type PokemonForm struct {
	ID               int64     `json:"id"`
	PokemonSpeciesID int64     `json:"pokemon_species_id"`
	FormKey          string    `json:"form_key"`
	NameJA           *string   `json:"name_ja"`
	NameEN           *string   `json:"name_en"`
	IsDefault        bool      `json:"is_default"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreatePokemonFormRequest struct {
	FormKey   string `json:"form_key"`
	NameJA    string `json:"name_ja"`
	NameEN    string `json:"name_en"`
	IsDefault *bool  `json:"is_default"`
	IsActive  *bool  `json:"is_active"`
}

func CreatePokemonForm(
	ctx context.Context,
	db *sql.DB,
	pokemonSpeciesID int64,
	request CreatePokemonFormRequest,
) (*PokemonForm, error) {
	request.FormKey = strings.TrimSpace(request.FormKey)
	request.NameJA = strings.TrimSpace(request.NameJA)
	request.NameEN = strings.TrimSpace(request.NameEN)

	if err := validateCreatePokemonFormRequest(
		pokemonSpeciesID,
		request,
	); err != nil {
		return nil, err
	}

	exists, err := pokemonSpeciesExists(ctx, db, pokemonSpeciesID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrPokemonSpeciesNotFound
	}

	isDefault := false
	if request.IsDefault != nil {
		isDefault = *request.IsDefault
	}

	isActive := true
	if request.IsActive != nil {
		isActive = *request.IsActive
	}

	const query = `
		INSERT INTO pokemon_forms (
			pokemon_species_id,
			form_key,
			name_ja,
			name_en,
			is_default,
			is_active
		)
		VALUES (
			$1,
			$2,
			NULLIF($3, ''),
			NULLIF($4, ''),
			$5,
			$6
		)
		RETURNING
			id,
			pokemon_species_id,
			form_key,
			name_ja,
			name_en,
			is_default,
			is_active,
			created_at,
			updated_at
	`

	var pokemonForm PokemonForm

	err = db.QueryRowContext(
		ctx,
		query,
		pokemonSpeciesID,
		request.FormKey,
		request.NameJA,
		request.NameEN,
		isDefault,
		isActive,
	).Scan(
		&pokemonForm.ID,
		&pokemonForm.PokemonSpeciesID,
		&pokemonForm.FormKey,
		&pokemonForm.NameJA,
		&pokemonForm.NameEN,
		&pokemonForm.IsDefault,
		&pokemonForm.IsActive,
		&pokemonForm.CreatedAt,
		&pokemonForm.UpdatedAt,
	)
	if err != nil {
		return nil, convertCreatePokemonFormError(err)
	}

	return &pokemonForm, nil
}

func FindPokemonFormsBySpeciesID(
	ctx context.Context,
	db *sql.DB,
	pokemonSpeciesID int64,
) ([]PokemonForm, error) {
	if pokemonSpeciesID <= 0 {
		return nil, ErrInvalidPokemonSpeciesID
	}

	exists, err := pokemonSpeciesExists(ctx, db, pokemonSpeciesID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrPokemonSpeciesNotFound
	}

	const query = `
		SELECT
			id,
			pokemon_species_id,
			form_key,
			name_ja,
			name_en,
			is_default,
			is_active,
			created_at,
			updated_at
		FROM pokemon_forms
		WHERE pokemon_species_id = $1
		ORDER BY
			is_default DESC,
			id ASC
	`

	rows, err := db.QueryContext(
		ctx,
		query,
		pokemonSpeciesID,
	)
	if err != nil {
		return nil, fmt.Errorf("find pokemon forms: %w", err)
	}
	defer rows.Close()

	pokemonForms := make([]PokemonForm, 0)

	for rows.Next() {
		var pokemonForm PokemonForm

		if err := rows.Scan(
			&pokemonForm.ID,
			&pokemonForm.PokemonSpeciesID,
			&pokemonForm.FormKey,
			&pokemonForm.NameJA,
			&pokemonForm.NameEN,
			&pokemonForm.IsDefault,
			&pokemonForm.IsActive,
			&pokemonForm.CreatedAt,
			&pokemonForm.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pokemon form: %w", err)
		}

		pokemonForms = append(pokemonForms, pokemonForm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pokemon forms: %w", err)
	}

	return pokemonForms, nil
}

func validateCreatePokemonFormRequest(
	pokemonSpeciesID int64,
	request CreatePokemonFormRequest,
) error {
	if pokemonSpeciesID <= 0 {
		return ErrInvalidPokemonSpeciesID
	}

	if request.FormKey == "" ||
		len(request.FormKey) > maxPokemonFormKeyLength ||
		!pokemonFormKeyPattern.MatchString(request.FormKey) {
		return ErrInvalidPokemonFormKey
	}

	if len([]rune(request.NameJA)) > maxPokemonFormNameLength {
		return ErrInvalidPokemonFormNameJA
	}

	if len([]rune(request.NameEN)) > maxPokemonFormNameLength {
		return ErrInvalidPokemonFormNameEN
	}

	return nil
}

func pokemonSpeciesExists(
	ctx context.Context,
	db *sql.DB,
	pokemonSpeciesID int64,
) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM pokemon_species
			WHERE id = $1
		)
	`

	var exists bool

	if err := db.QueryRowContext(
		ctx,
		query,
		pokemonSpeciesID,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf(
			"check pokemon species existence: %w",
			err,
		)
	}

	return exists, nil
}

func convertCreatePokemonFormError(err error) error {
	var pgError *pgconn.PgError
	if !errors.As(err, &pgError) {
		return fmt.Errorf("create pokemon form: %w", err)
	}

	switch pgError.Code {
	case "23503":
		return ErrPokemonSpeciesNotFound

	case "23505":
		switch pgError.ConstraintName {
		case "uq_pokemon_forms_species_form_key":
			return ErrDuplicatePokemonForm

		case "uq_pokemon_forms_default":
			return ErrDuplicateDefaultForm
		}
	}

	return fmt.Errorf("create pokemon form: %w", err)
}
