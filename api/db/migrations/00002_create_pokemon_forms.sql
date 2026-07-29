-- +goose Up
CREATE TABLE pokemon_forms (
    id BIGSERIAL PRIMARY KEY,
    pokemon_species_id BIGINT NOT NULL,
    form_key VARCHAR(100) NOT NULL,
    name_ja VARCHAR(100),
    name_en VARCHAR(100),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pokemon_forms_species
        FOREIGN KEY (pokemon_species_id)
        REFERENCES pokemon_species (id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_pokemon_forms_species_form_key
        UNIQUE (pokemon_species_id, form_key),

    CONSTRAINT chk_pokemon_forms_form_key
        CHECK (form_key ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$')
);

CREATE UNIQUE INDEX uq_pokemon_forms_default
    ON pokemon_forms (pokemon_species_id)
    WHERE is_default = TRUE;

CREATE INDEX idx_pokemon_forms_species_id
    ON pokemon_forms (pokemon_species_id);

-- +goose Down
DROP TABLE IF EXISTS pokemon_forms;
