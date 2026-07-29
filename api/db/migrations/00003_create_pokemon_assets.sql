-- +goose Up
CREATE TABLE pokemon_assets (
    id BIGSERIAL PRIMARY KEY,
    pokemon_form_id BIGINT NOT NULL,
    asset_type VARCHAR(50) NOT NULL,
    storage_path VARCHAR(500) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    width INTEGER,
    height INTEGER,
    frame_count INTEGER NOT NULL DEFAULT 1,
    frame_width INTEGER,
    frame_height INTEGER,
    frame_rate NUMERIC(6, 2),
    is_loop BOOLEAN NOT NULL DEFAULT FALSE,
    checksum_sha256 CHAR(64),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pokemon_assets_form
        FOREIGN KEY (pokemon_form_id)
        REFERENCES pokemon_forms (id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_pokemon_assets_storage_path
        UNIQUE (storage_path),

    CONSTRAINT chk_pokemon_assets_asset_type
        CHECK (
            asset_type IN (
                'image',
                'idle',
                'enter',
                'physical_attack',
                'special_attack',
                'damage',
                'faint',
                'victory'
            )
        ),

    CONSTRAINT chk_pokemon_assets_file_size
        CHECK (file_size > 0),

    CONSTRAINT chk_pokemon_assets_image_size
        CHECK (
            (width IS NULL OR width > 0)
            AND (height IS NULL OR height > 0)
        ),

    CONSTRAINT chk_pokemon_assets_frame
        CHECK (
            frame_count > 0
            AND (frame_width IS NULL OR frame_width > 0)
            AND (frame_height IS NULL OR frame_height > 0)
            AND (frame_rate IS NULL OR frame_rate > 0)
        ),

    CONSTRAINT chk_pokemon_assets_checksum
        CHECK (
            checksum_sha256 IS NULL
            OR checksum_sha256 ~ '^[0-9a-f]{64}$'
        )
);

CREATE INDEX idx_pokemon_assets_form_id
    ON pokemon_assets (pokemon_form_id);

CREATE INDEX idx_pokemon_assets_form_type
    ON pokemon_assets (pokemon_form_id, asset_type);

-- +goose Down
DROP TABLE IF EXISTS pokemon_assets;
