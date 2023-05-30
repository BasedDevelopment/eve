-- +goose Up
-- +goose StatementBegin
INSERT INTO public.hv (id, hostname, auto_url, auto_serial, site) VALUES ('85833bb8-2f0a-4b1e-981f-f9cb3597904c', 'hv0.sit.bns.sh', 'hv0.sit.bns.sh:3001','1683180991', 'sit');
INSERT INTO public.profile (id, name, email, is_admin, password) VALUES ('1636cad3-f638-4bb3-b0f2-dbe5fafe9b6e', 'Eric', 'admin@ericz.me', TRUE, '$2a$10$RwxoIEyBvbDutg6vYzt0ceiBEyqjzHI/21r4vOZwi0afQqe0LzY/6');
INSERT INTO public.profile (id, name, email, is_admin, password) VALUES ('b7549879-700d-4ee9-abb2-fe438e7eb133', 'Eric', 'user@ericz.me', FALSE, '$2a$10$RwxoIEyBvbDutg6vYzt0ceiBEyqjzHI/21r4vOZwi0afQqe0LzY/6');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE sessions;
TRUNCATE vm CASCADE;
TRUNCATE hv CASCADE;
TRUNCATE profile CASCADE;
-- +goose StatementEnd
