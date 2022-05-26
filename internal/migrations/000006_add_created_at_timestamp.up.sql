BEGIN;

ALTER TABLE
    syllabuses
ADD
    COLUMN created_at timestamp NOT NULL;

ALTER TABLE
    attachments
ADD
    COLUMN created_at timestamp NOT NULL;

ALTER TABLE
    magic_tokens
ADD
    COLUMN created_at timestamp NOT NULL;

ALTER TABLE
    contributors
ADD
    COLUMN created_at timestamp NOT NULL;

COMMIT;