-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.profile (
    id uuid NOT NULL PRIMARY KEY,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL UNIQUE,
    password text NOT NULL,
    disabled boolean NOT NULL DEFAULT FALSE,
    is_admin boolean NOT NULL DEFAULT FALSE,
    last_login timestamp with time zone NOT NULL DEFAULT now(),
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text NOT NULL DEFAULT ''
);

CREATE TABLE public.sessions (
    owner uuid NOT NULL REFERENCES public.profile(id),
    token_public character varying(255) NOT NULL PRIMARY KEY,
    token_secret character varying(255) NOT NULL,
    token_salt character varying(255) NOT NULL,
    token_version character varying(4) NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    expires timestamp with time zone NOT NULL
);

CREATE TABLE public.hv (
    id uuid NOT NULL PRIMARY KEY,
    hostname character varying(255) NOT NULL,
    ip inet NOT NULL,
    port integer NOT NULL DEFAULT 16509,
    site character varying(255) NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text NOT NULL DEFAULT ''
);

CREATE TABLE public.vm (
    id uuid NOT NULL PRIMARY KEY,
    hv_id uuid NOT NULL REFERENCES hv (id),
    hostname character varying(255) NOT NULL,
    profile_id uuid NOT NULL REFERENCES profile (id),
    cpu integer NOT NULL,
    memory bigint NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text NOT NULL DEFAULT ''
);

CREATE TABLE public.vm_nic (
    id uuid NOT NULL PRIMARY KEY,
    vm_id uuid NOT NULL REFERENCES vm (id),
    name character varying(255) NOT NULL,
    mac macaddr NOT NULL,
    ips inet[] NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text NOT NULL DEFAULT ''
);

CREATE TABLE public.vm_storage (
    id uuid NOT NULL PRIMARY KEY,
    vm_id uuid NOT NULL REFERENCES vm (id),
    size bigint NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text NOT NULL DEFAULT ''
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.vm_storage;
DROP TABLE public.vm_nic;
DROP TABLE public.vm;
DROP TABLE public.hv;
DROP TABLE public.sessions;
DROP TABLE public.profile;
-- +goose StatementEnd
