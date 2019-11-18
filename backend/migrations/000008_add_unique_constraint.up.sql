BEGIN;

CREATE UNIQUE INDEX subscription_user_state_unique_idx ON subscription_user_state (subscription_id, user_id);

ALTER TABLE subscription_user_state ADD CONSTRAINT subscription_user_state_unique UNIQUE USING INDEX subscription_user_state_unique_idx;

COMMIT;
