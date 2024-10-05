-- +goose Up
-- +goose StatementBegin
CREATE TABLE projects(
    id integer primary key autoincrement,
    parent_id integer nullable references projects(id) ON DELETE CASCADE,
    title varchar(255) not null,
    description varchar(255) not null,
    status varchar(255) not null,
    assignee_id integer nullable references users(id),
    urgency tinyint not null,
    notes varchar(255),
    start_date timestamp nullable,
    end_date timestamp nullable,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (title)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE projects;
-- +goose StatementEnd
