-- +goose Up

CREATE SCHEMA IF NOT EXISTS app;

-- ---------------------------------------------------------------------------

-- +goose StatementBegin

CREATE
    OR REPLACE FUNCTION set_updated_at()
    RETURNS trigger AS
$$
BEGIN
    NEW.updated_at
        = NOW();
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

-- +goose StatementEnd

-- ---------------------------------------------------------------------------

CREATE TABLE app.skills
(
    id                 UUID        NOT NULL DEFAULT gen_random_uuid(),
    CONSTRAINT pk_skills PRIMARY KEY (id),

    name               TEXT        NOT NULL,
    CONSTRAINT uq_skills_name UNIQUE (name),

    youtube_video_id   TEXT        NOT NULL,
    CONSTRAINT uq_skills_youtube_video_id UNIQUE (youtube_video_id),

    is_video_landscape BOOLEAN     NOT NULL DEFAULT FALSE,

    level              SMALLINT    NOT NULL,
    CONSTRAINT ck_skills_level CHECK (level >= 1 AND level <= 8),

    categories         TEXT[]      NOT NULL DEFAULT '{}',
    CONSTRAINT ck_skills_categories CHECK
        (categories <@ ARRAY ['basics', 'manipulation', 'footwork', 'backward', 'wraps', 'releases', 'floaters', 'multiples']::TEXT[]),

    prerequisites      UUID[]      NOT NULL DEFAULT '{}',
    CONSTRAINT ck_skills_prerequisites CHECK (NOT (id = ANY (prerequisites))),

    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_skills_updated_at
    BEFORE UPDATE
    ON app.skills
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();