CREATE TABLE IF NOT EXISTS magic_tokens(
 id SERIAL PRIMARY KEY,
 token BYTEA NOT NULL,
 syllabus_token_id BIGINT NOT NUll 
);