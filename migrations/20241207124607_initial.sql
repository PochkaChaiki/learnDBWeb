-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS user(
    user_id  integer primary key,
    username text not null,
    password text not null
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS db(
    db_id   integer primary key,
    db_name text not null
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS db_sample(
    db_sample_id integer primary key,
    description  text,
    filepath     text not null,
    db_id        integer REFERENCES db(db_id) ON UPDATE CASCADE ON DELETE SET NULL
);

-- +goose StatementEnd
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS question(
    question_id    integer primary key,
    question_text  text not null,
    correct_answer text not null,
    db_sample_id   integer REFERENCES db_sample(db_sample_id) ON UPDATE CASCADE ON DELETE SET NULL
);

-- +goose StatementEnd
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS query(
    query_id    integer primary key,
    script      text not null,
    info        text default null,
    executed_at text default (datetime('now')),
    user_id     integer REFERENCES user(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    db_id       integer REFERENCES db(db_id) ON UPDATE CASCADE ON DELETE SET NULL
);

-- +goose StatementEnd
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS answer(
    answer_id   integer primary key,
    answer_text text,
    is_correct  integer,
    question_id integer REFERENCES question(question_id) ON UPDATE CASCADE ON DELETE SET NULL,
    query_id    integer REFERENCES query(query_id) ON UPDATE CASCADE ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
DROP TABLE db;
DROP TABLE db_sample;
DROP TABLE question;
DROP TABLE query;
DROP TABLE answer;

-- +goose StatementEnd
