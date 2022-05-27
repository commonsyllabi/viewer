BEGIN;

ALTER TABLE
    syllabuses DROP COLUMN created_at;

ALTER TABLE
    attachments DROP COLUMN created_at;

ALTER TABLE
    magic_tokens DROP COLUMN created_at;

ALTER TABLE
    contributors DROP COLUMN created_at;

COMMIT;