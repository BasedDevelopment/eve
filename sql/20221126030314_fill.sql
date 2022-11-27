-- +goose Up
-- +goose StatementBegin
INSERT INTO public.hv (id, hostname, ip, port, site) VALUES ('85833bb8-2f0a-4b1e-981f-f9cb3597904c', 'dorm0.sit.eric.si', '10.10.9.4', 16509, 'sit');
INSERT INTO public.profile (id, name, email, is_admin, password) VALUES ('1636cad3-f638-4bb3-b0f2-dbe5fafe9b6e', 'Eric', 'eric@ericz.me', TRUE, '$2a$10$RwxoIEyBvbDutg6vYzt0ceiBEyqjzHI/21r4vOZwi0afQqe0LzY/6');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.hv WHERE id='85833bb8-2f0a-4b1e-981f-f9cb3597904c';
DELETE FROM public.profile WHERE id='1636cad3-f638-4bb3-b0f2-dbe5fafe9b6e';
-- +goose StatementEnd
