-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.profile (
    id uuid NOT NULL PRIMARY KEY,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password text NOT NULL,
    disabled boolean NOT NULL DEFAULT FALSE,
    is_admin boolean NOT NULL DEFAULT FALSE,
    last_login timestamp with time zone,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text
);

CREATE TABLE public.token (
    token_public character varying(255) NOT NULL PRIMARY KEY,
    token_private character varying(255) NOT NULL,
    profile_id uuid NOT NULL REFERENCES public.profile(id),
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
    remarks text
);

CREATE TABLE public.hv_nic (
    id uuid NOT NULL PRIMARY KEY,
    hv_id uuid NOT NULL REFERENCES public.hv(id),
    name character varying(255) NOT NULL,
    mac character varying(255) NOT NULL,
    ip inet[] NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text
);

CREATE TABLE public.hv_storage (
    id uuid NOT NULL PRIMARY KEY,
    hv_id uuid NOT NULL REFERENCES hv (id),
    size integer NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text
);

CREATE TABLE public.vm (
    id uuid NOT NULL PRIMARY KEY,
    hv_id uuid NOT NULL REFERENCES hv (id),
    hostname character varying(255) NOT NULL,
    profile_id uuid NOT NULL REFERENCES profile (id),
    cpu integer NOT NULL,
    memory integer NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text
);

CREATE TABLE public.vm_nic (
    id uuid NOT NULL PRIMARY KEY,
    vm_id uuid NOT NULL REFERENCES vm (id),
    name character varying(255) NOT NULL,
    mac macaddr NOT NULL,
    ips inet[] NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text
);

CREATE TABLE public.vm_storage (
    id uuid NOT NULL PRIMARY KEY,
    vm_id uuid NOT NULL REFERENCES vm (id),
    size bigint NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    updated timestamp with time zone NOT NULL DEFAULT now(),
    remarks text
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.vm_storage;
DROP TABLE public.vm_nic;
DROP TABLE public.vm;
DROP TABLE public.hv_nic;
DROP TABLE public.hv_storage;
DROP TABLE public.hv;
DROP TABLE public.token;
DROP TABLE public.profile;
-- +goose StatementEnd
