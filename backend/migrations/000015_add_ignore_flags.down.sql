BEGIN;

ALTER TABLE subscription DROP COLUMN ignore_rt;

ALTER TABLE subscription DROP COLUMN ignore_replies;

COMMIT;
