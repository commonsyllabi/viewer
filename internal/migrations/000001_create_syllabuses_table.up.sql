CREATE TABLE IF NOT EXISTS syllabuses(
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description VARCHAR,
    email VARCHAR,
    contributor_id BIGINT
);