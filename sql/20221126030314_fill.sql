-- +goose Up
-- +goose StatementBegin
INSERT INTO public.hv (id, hostname, ip, port, site) VALUES ('85833bb8-2f0a-4b1e-981f-f9cb3597904c', 'dorm0.sit.eric.si', '10.10.9.4', 16509, 'sit');
INSERT INTO public.hv (id, hostname, ip, port, site, remarks) VALUES ('988d4d3d-289d-4d70-8595-265fd9e1cab5', 'broke0.sit.eric.si', '192.168.0.1', 16509, 'sit', 'brok');
INSERT INTO public.profile (id, name, email, is_admin, password) VALUES ('1636cad3-f638-4bb3-b0f2-dbe5fafe9b6e', 'Eric', 'admin@ericz.me', TRUE, '$2a$10$RwxoIEyBvbDutg6vYzt0ceiBEyqjzHI/21r4vOZwi0afQqe0LzY/6');
INSERT INTO public.profile (id, name, email, is_admin, password) VALUES ('b7549879-700d-4ee9-abb2-fe438e7eb133', 'Eric', 'user@ericz.me', FALSE, '$2a$10$RwxoIEyBvbDutg6vYzt0ceiBEyqjzHI/21r4vOZwi0afQqe0LzY/6');
INSERT INTO public.vm (id, hv_id, hostname, profile_id, cpu, memory) VALUES ('7f119176-4a59-4ce9-adbc-433011e5b5bb', '85833bb8-2f0a-4b1e-981f-f9cb3597904c', 'debtest.eric.si', '1636cad3-f638-4bb3-b0f2-dbe5fafe9b6e', '2', '2048');
INSERT INTO public.vm (id, hv_id, hostname, profile_id, cpu, memory) VALUES ('fb667190-6967-4d2b-8ffb-2838ee445f2b', '85833bb8-2f0a-4b1e-981f-f9cb3597904c', 'debtest2.eric.si', 'b7549879-700d-4ee9-abb2-fe438e7eb133', '2', '2048');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE sessions;
TRUNCATE vm CASCADE;
TRUNCATE hv CASCADE;
TRUNCATE profile CASCADE;
-- +goose StatementEnd
