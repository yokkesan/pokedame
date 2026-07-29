-- +goose Up
CREATE TABLE pokemon_species (
    id BIGSERIAL PRIMARY KEY,
    national_dex_number INTEGER NOT NULL,
    name_ja VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    slug VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_pokemon_species_national_dex_number
        UNIQUE (national_dex_number),

    CONSTRAINT uq_pokemon_species_slug
        UNIQUE (slug),

    CONSTRAINT chk_pokemon_species_national_dex_number
        CHECK (national_dex_number > 0),

    CONSTRAINT chk_pokemon_species_slug
        CHECK (slug ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$')
);

-- +goose Down
DROP TABLE IF EXISTS pokemon_species;
