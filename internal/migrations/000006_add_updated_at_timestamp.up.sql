BEGIN;

ALTER TABLE
    syllabuses
ADD
    COLUMN updated_at timestamp NOT NULL;

ALTER TABLE
    attachments
ADD
    COLUMN updated_at timestamp NOT NULL;

ALTER TABLE
    magic_tokens
ADD
    COLUMN updated_at timestamp NOT NULL;

ALTER TABLE
    contributors
ADD
    COLUMN updated_at timestamp NOT NULL;

COMMIT;