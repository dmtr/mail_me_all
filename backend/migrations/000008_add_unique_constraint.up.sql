BEGIN;

CREATE UNIQUE INDEX subscription_user_state_unique_idx ON subscription_user_state (subscription_id, user_id);

ALTER TABLE subscription_user_state ADD CONSTRAINT subscription_user_state_unique UNIQUE USING INDEX subscription_user_state_unique_idx;

ALTER TABLE subscription_user_state DROP COLUMN user_id;

ALTER TABLE subscription_user_state ADD COLUMN user_twitter_id VARCHAR NOT NULL;

ALTER TABLE subscription_user_state ADD CONSTRAINT subscription_user_state_user_id_fk FOREIGN KEY (user_twitter_id) REFERENCES subscription_user (twitter_id) ON DELETE CASCADE;

COMMIT;
