BEGIN;

ALTER TABLE subscription ADD COLUMN ignore_rt BOOLEAN DEFAULT false;

ALTER TABLE subscription ADD COLUMN ignore_replies BOOLEAN DEFAULT false;

COMMIT;
