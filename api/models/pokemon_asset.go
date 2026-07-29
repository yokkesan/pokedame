package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidPokemonFormID      = errors.New("invalid pokemon form id")
	ErrPokemonFormNotFound       = errors.New("pokemon form not found")
	ErrInvalidPokemonAssetID     = errors.New("invalid pokemon asset id")
	ErrPokemonAssetNotFound      = errors.New("pokemon asset not found")
	ErrInvalidPokemonAssetType   = errors.New("invalid pokemon asset type")
	ErrDuplicateAssetStoragePath = errors.New("duplicate asset storage path")
)

var validPokemonAssetTypes = map[string]struct{}{
	"image":           {},
	"idle":            {},
	"enter":           {},
	"physical_attack": {},
	"special_attack":  {},
	"damage":          {},
	"faint":           {},
	"victory":         {},
}

type PokemonAsset struct {
	ID               int64     `json:"id"`
	PokemonFormID    int64     `json:"pokemon_form_id"`
	AssetType        string    `json:"asset_type"`
	StoragePath      string    `json:"storage_path"`
	OriginalFilename string    `json:"original_filename"`
	MimeType         string    `json:"mime_type"`
	FileSize         int64     `json:"file_size"`
	Width            *int      `json:"width"`
	Height           *int      `json:"height"`
	FrameCount       int       `json:"frame_count"`
	FrameWidth       *int      `json:"frame_width"`
	FrameHeight      *int      `json:"frame_height"`
	FrameRate        *float64  `json:"frame_rate"`
	IsLoop           bool      `json:"is_loop"`
	ChecksumSHA256   *string   `json:"checksum_sha256"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreatePokemonAssetRequest struct {
	PokemonFormID    int64
	AssetType        string
	StoragePath      string
	OriginalFilename string
	MimeType         string
	FileSize         int64
	Width            *int
	Height           *int
	FrameCount       int
	FrameWidth       *int
	FrameHeight      *int
	FrameRate        *float64
	IsLoop           bool
	ChecksumSHA256   string
	IsActive         bool
}

func CreatePokemonAsset(
	ctx context.Context,
	db *sql.DB,
	request CreatePokemonAssetRequest,
) (*PokemonAsset, error) {
	if request.PokemonFormID <= 0 {
		return nil, ErrInvalidPokemonFormID
	}

	if _, ok := validPokemonAssetTypes[request.AssetType]; !ok {
		return nil, ErrInvalidPokemonAssetType
	}

	exists, err := pokemonFormExists(ctx, db, request.PokemonFormID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrPokemonFormNotFound
	}

	const query = `
		INSERT INTO pokemon_assets (
			pokemon_form_id,
			asset_type,
			storage_path,
			original_filename,
			mime_type,
			file_size,
			width,
			height,
			frame_count,
			frame_width,
			frame_height,
			frame_rate,
			is_loop,
			checksum_sha256,
			is_active
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, NULLIF($14, ''), $15
		)
		RETURNING
			id,
			pokemon_form_id,
			asset_type,
			storage_path,
			original_filename,
			mime_type,
			file_size,
			width,
			height,
			frame_count,
			frame_width,
			frame_height,
			frame_rate,
			is_loop,
			checksum_sha256,
			is_active,
			created_at,
			updated_at
	`

	var asset PokemonAsset

	err = db.QueryRowContext(
		ctx,
		query,
		request.PokemonFormID,
		request.AssetType,
		request.StoragePath,
		request.OriginalFilename,
		request.MimeType,
		request.FileSize,
		request.Width,
		request.Height,
		request.FrameCount,
		request.FrameWidth,
		request.FrameHeight,
		request.FrameRate,
		request.IsLoop,
		request.ChecksumSHA256,
		request.IsActive,
	).Scan(
		&asset.ID,
		&asset.PokemonFormID,
		&asset.AssetType,
		&asset.StoragePath,
		&asset.OriginalFilename,
		&asset.MimeType,
		&asset.FileSize,
		&asset.Width,
		&asset.Height,
		&asset.FrameCount,
		&asset.FrameWidth,
		&asset.FrameHeight,
		&asset.FrameRate,
		&asset.IsLoop,
		&asset.ChecksumSHA256,
		&asset.IsActive,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23503":
				return nil, ErrPokemonFormNotFound

			case "23505":
				if pgError.ConstraintName == "uq_pokemon_assets_storage_path" {
					return nil, ErrDuplicateAssetStoragePath
				}
			}
		}

		return nil, fmt.Errorf("create pokemon asset: %w", err)
	}

	return &asset, nil
}

func FindPokemonAssetsByFormID(
	ctx context.Context,
	db *sql.DB,
	pokemonFormID int64,
) ([]PokemonAsset, error) {
	if pokemonFormID <= 0 {
		return nil, ErrInvalidPokemonFormID
	}

	exists, err := pokemonFormExists(ctx, db, pokemonFormID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrPokemonFormNotFound
	}

	const query = `
		SELECT
			id,
			pokemon_form_id,
			asset_type,
			storage_path,
			original_filename,
			mime_type,
			file_size,
			width,
			height,
			frame_count,
			frame_width,
			frame_height,
			frame_rate,
			is_loop,
			checksum_sha256,
			is_active,
			created_at,
			updated_at
		FROM pokemon_assets
		WHERE pokemon_form_id = $1
		ORDER BY id ASC
	`

	rows, err := db.QueryContext(ctx, query, pokemonFormID)
	if err != nil {
		return nil, fmt.Errorf("find pokemon assets: %w", err)
	}
	defer rows.Close()

	assets := make([]PokemonAsset, 0)

	for rows.Next() {
		var asset PokemonAsset

		if err := rows.Scan(
			&asset.ID,
			&asset.PokemonFormID,
			&asset.AssetType,
			&asset.StoragePath,
			&asset.OriginalFilename,
			&asset.MimeType,
			&asset.FileSize,
			&asset.Width,
			&asset.Height,
			&asset.FrameCount,
			&asset.FrameWidth,
			&asset.FrameHeight,
			&asset.FrameRate,
			&asset.IsLoop,
			&asset.ChecksumSHA256,
			&asset.IsActive,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pokemon asset: %w", err)
		}

		assets = append(assets, asset)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pokemon assets: %w", err)
	}

	return assets, nil
}

func FindPokemonAssetByID(
	ctx context.Context,
	db *sql.DB,
	assetID int64,
) (*PokemonAsset, error) {
	if assetID <= 0 {
		return nil, ErrInvalidPokemonAssetID
	}

	const query = `
		SELECT
			id,
			pokemon_form_id,
			asset_type,
			storage_path,
			original_filename,
			mime_type,
			file_size,
			width,
			height,
			frame_count,
			frame_width,
			frame_height,
			frame_rate,
			is_loop,
			checksum_sha256,
			is_active,
			created_at,
			updated_at
		FROM pokemon_assets
		WHERE id = $1
	`

	var asset PokemonAsset

	err := db.QueryRowContext(ctx, query, assetID).Scan(
		&asset.ID,
		&asset.PokemonFormID,
		&asset.AssetType,
		&asset.StoragePath,
		&asset.OriginalFilename,
		&asset.MimeType,
		&asset.FileSize,
		&asset.Width,
		&asset.Height,
		&asset.FrameCount,
		&asset.FrameWidth,
		&asset.FrameHeight,
		&asset.FrameRate,
		&asset.IsLoop,
		&asset.ChecksumSHA256,
		&asset.IsActive,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPokemonAssetNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find pokemon asset: %w", err)
	}

	return &asset, nil
}

func DeletePokemonAsset(
	ctx context.Context,
	db *sql.DB,
	assetID int64,
) error {
	if assetID <= 0 {
		return ErrInvalidPokemonAssetID
	}

	const query = `
		DELETE FROM pokemon_assets
		WHERE id = $1
	`

	result, err := db.ExecContext(ctx, query, assetID)
	if err != nil {
		return fmt.Errorf("delete pokemon asset: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted pokemon asset rows: %w", err)
	}

	if affectedRows == 0 {
		return ErrPokemonAssetNotFound
	}

	return nil
}

func pokemonFormExists(
	ctx context.Context,
	db *sql.DB,
	pokemonFormID int64,
) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM pokemon_forms
			WHERE id = $1
		)
	`

	var exists bool

	if err := db.QueryRowContext(
		ctx,
		query,
		pokemonFormID,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf(
			"check pokemon form existence: %w",
			err,
		)
	}

	return exists, nil
}
