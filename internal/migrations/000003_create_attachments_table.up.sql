CREATE TABLE IF NOT EXISTS attachments(
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    file BYTEA NOT NULL,
    type VARCHAR NOT NULL,
    syllabus_attached_id BIGINT
);