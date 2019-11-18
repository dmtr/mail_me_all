BEGIN;

ALTER TABLE subscription_user_state DROP CONSTRAINT subscription_user_state_unique;

ALTER TABLE subscription_user_state DROP COLUMN user_twitter_id;

ALTER TABLE subscription_user_state ADD COLUMN user_id INTEGER NOT NULL;

ALTER TABLE subscription_user_state ADD CONSTRAINT subscription_user_state_user_id_fk FOREIGN KEY (user_id) REFERENCES subscription_user (id) ON DELETE CASCADE;

COMMIT;
