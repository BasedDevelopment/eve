-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS public.token;

CREATE TABLE public.sessions (
    owner uuid NOT NULL REFERENCES public.profile(id),
    token_public character varying(255) NOT NULL PRIMARY KEY,
    token_secret character varying(255) NOT NULL,
    token_salt character varying(255) NOT NULL,
    token_version character varying(4) NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    expires timestamp with time zone NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.sessions;

CREATE TABLE IF NOT EXISTS public.token (
    token_public character varying(255) NOT NULL PRIMARY KEY,
    token_private character varying(255) NOT NULL,
    profile_id uuid NOT NULL REFERENCES public.profile(id),
    created timestamp with time zone NOT NULL DEFAULT now(),
    expires timestamp with time zone NOT NULL
);

-- +goose StatementEnd
