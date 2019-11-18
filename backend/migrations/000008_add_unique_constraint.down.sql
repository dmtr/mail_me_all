BEGIN;

ALTER TABLE subscription_user_state DROP CONSTRAINT subscription_user_state_unique;

COMMIT;
