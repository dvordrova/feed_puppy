-- +goose Up
-- +goose StatementBegin
CREATE TABLE dogs (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL,
    birth_date TEXT    NOT NULL,
    sex        TEXT    NOT NULL
                       CHECK(sex IN ('Male', 'Female')),
    breed      TEXT    NOT NULL
);

CREATE TABLE users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id INTEGER UNIQUE NOT NULL,
    name        TEXT    NOT NULL,
    language    TEXT    NOT NULL
                        CHECK(language IN ('en', 'ru')),
    current_dog INTEGER,
    state       TEXT    NOT NULL
                        CHECK(state IN ('new_dog',
                                        'action_feed',
                                        'action_yum',
                                        'action_weigh',
                                        'dog_selected',
                                        'change_name'))
                        DEFAULT 'new_dog',
    FOREIGN KEY (current_dog) REFERENCES dogs(id)
);

CREATE TABLE foods (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    dog_id           INTEGER NOT NULL,
    title            TEXT    NOT NULL,
    calories         INTEGER NOT NULL,
    is_last_selected BOOLEAN NOT NULL,
    FOREIGN KEY (dog_id) REFERENCES dogs(id)
);

CREATE TABLE action_descriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL
);

INSERT INTO action_descriptions (type) VALUES
    ('action_feed'),
    ('action_yum'),
    ('action_vacine'),
    ('action_weigh'),
    ('action_shit_detected'),
    ('action_shit_removed'),
    ('action_pee_detected'),
    ('action_pee_removed'),
    ('action_play_start'),
    ('action_play_end'),
    ('action_teach_start'),
    ('action_teach_end'),
    ('action_sleep_start'),
    ('action_sleep_end'),
    ('action_walk_start'),
    ('action_walk_end');

CREATE TABLE actions (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id          INTEGER NOT NULL,
    dog_id           INTEGER NOT NULL,
    timestamp        INTEGER NOT NULL,
    action_id        INTEGER NOT NULL,
    additional_info  TEXT,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (dog_id) REFERENCES dogs(id),
    FOREIGN KEY (action_id) REFERENCES actions(id)
);

CREATE TABLE subscriptions (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    dog_id  INTEGER NOT NULL,
    type    TEXT    NOT NULL
                    CHECK(type IN ('owner', 'reader')),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (dog_id) REFERENCES dogs(id)
);

CREATE TABLE invite_subscriptions (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    hash    TEXT    NOT NULL,
    dog_id  INTEGER NOT NULL,
    type    TEXT    NOT NULL
                    CHECK(type IN ('owner', 'reader')),
    used    BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose StatementEnd
