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
	maxPokemonNameLength = 100
	maxPokemonSlugLength = 100
)

var (
	ErrInvalidNationalDexNumber = errors.New("invalid national dex number")
	ErrInvalidPokemonNameJA     = errors.New("invalid pokemon japanese name")
	ErrInvalidPokemonNameEN     = errors.New("invalid pokemon english name")
	ErrInvalidPokemonSlug       = errors.New("invalid pokemon slug")
	ErrDuplicatePokemonSpecies  = errors.New("duplicate pokemon species")
)

var pokemonSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type PokemonSpecies struct {
	ID                int64     `json:"id"`
	NationalDexNumber int       `json:"national_dex_number"`
	NameJA            string    `json:"name_ja"`
	NameEN            *string   `json:"name_en"`
	Slug              string    `json:"slug"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreatePokemonSpeciesRequest struct {
	NationalDexNumber int    `json:"national_dex_number"`
	NameJA            string `json:"name_ja"`
	NameEN            string `json:"name_en"`
	Slug              string `json:"slug"`
	IsActive          *bool  `json:"is_active"`
}

func CreatePokemonSpecies(
	ctx context.Context,
	db *sql.DB,
	request CreatePokemonSpeciesRequest,
) (*PokemonSpecies, error) {
	request.NameJA = strings.TrimSpace(request.NameJA)
	request.NameEN = strings.TrimSpace(request.NameEN)
	request.Slug = strings.TrimSpace(request.Slug)

	if err := validateCreatePokemonSpeciesRequest(request); err != nil {
		return nil, err
	}

	isActive := true
	if request.IsActive != nil {
		isActive = *request.IsActive
	}

	const query = `
		INSERT INTO pokemon_species (
			national_dex_number,
			name_ja,
			name_en,
			slug,
			is_active
		)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5)
		RETURNING
			id,
			national_dex_number,
			name_ja,
			name_en,
			slug,
			is_active,
			created_at,
			updated_at
	`

	var pokemon PokemonSpecies

	err := db.QueryRowContext(
		ctx,
		query,
		request.NationalDexNumber,
		request.NameJA,
		request.NameEN,
		request.Slug,
		isActive,
	).Scan(
		&pokemon.ID,
		&pokemon.NationalDexNumber,
		&pokemon.NameJA,
		&pokemon.NameEN,
		&pokemon.Slug,
		&pokemon.IsActive,
		&pokemon.CreatedAt,
		&pokemon.UpdatedAt,
	)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return nil, ErrDuplicatePokemonSpecies
		}

		return nil, fmt.Errorf("create pokemon species: %w", err)
	}

	return &pokemon, nil
}

func FindAllPokemonSpecies(
	ctx context.Context,
	db *sql.DB,
) ([]PokemonSpecies, error) {
	const query = `
		SELECT
			id,
			national_dex_number,
			name_ja,
			name_en,
			slug,
			is_active,
			created_at,
			updated_at
		FROM pokemon_species
		ORDER BY national_dex_number ASC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find pokemon species: %w", err)
	}
	defer rows.Close()

	pokemonList := make([]PokemonSpecies, 0)

	for rows.Next() {
		var pokemon PokemonSpecies

		if err := rows.Scan(
			&pokemon.ID,
			&pokemon.NationalDexNumber,
			&pokemon.NameJA,
			&pokemon.NameEN,
			&pokemon.Slug,
			&pokemon.IsActive,
			&pokemon.CreatedAt,
			&pokemon.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pokemon species: %w", err)
		}

		pokemonList = append(pokemonList, pokemon)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pokemon species: %w", err)
	}

	return pokemonList, nil
}

func validateCreatePokemonSpeciesRequest(
	request CreatePokemonSpeciesRequest,
) error {
	if request.NationalDexNumber <= 0 {
		return ErrInvalidNationalDexNumber
	}

	if request.NameJA == "" ||
		len([]rune(request.NameJA)) > maxPokemonNameLength {
		return ErrInvalidPokemonNameJA
	}

	if len([]rune(request.NameEN)) > maxPokemonNameLength {
		return ErrInvalidPokemonNameEN
	}

	if request.Slug == "" ||
		len(request.Slug) > maxPokemonSlugLength ||
		!pokemonSlugPattern.MatchString(request.Slug) {
		return ErrInvalidPokemonSlug
	}

	return nil
}
