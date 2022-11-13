CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE notes (
    id INTEGER PRIMARY KEY,
    header VARCHAR(255) NOT NULL,
    short_body VARCHAR(255),
    color VARCHAR(6) NOT NULL,
    edited DATETIME
);

CREATE TABLE notes_body (
    id INTEGER REFERENCES notes(id) ON DELETE CASCADE NOT NULL,
    body TEXT
);

CREATE TABLE users_notes (
    id INTEGER PRIMARY KEY,
    users_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    notes_id INTEGER REFERENCES notes(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(6) NOT NULL
);

CREATE TABLE tags_notes (
    id INTEGER PRIMARY KEY,
    tags_id INT REFERENCES tags(id) ON DELETE CASCADE NOT NULL,
    notes_id INT REFERENCES notes(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE users_tags (
    id INTEGER PRIMARY KEY,
    users_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    tags_id INT REFERENCES tags(id) ON DELETE CASCADE NOT NULL
);