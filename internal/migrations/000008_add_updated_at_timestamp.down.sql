BEGIN;

ALTER TABLE
    syllabuses DROP COLUMN updated_at;

ALTER TABLE
    attachments DROP COLUMN updated_at;

ALTER TABLE
    magic_tokens DROP COLUMN updated_at;

ALTER TABLE
    contributors DROP COLUMN updated_at;

COMMIT;