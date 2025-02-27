-- +goose Up
-- +goose StatementBegin
CREATE TABLE if NOT EXISTS tasks (
    id            serial primary key,
    title       text     not null,
    description        int      ,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
