BEGIN;

ALTER TABLE subscription_user_state DROP COLUMN user_id;

ALTER TABLE subscription_user_state ADD COLUMN user_twitter_id VARCHAR NOT NULL;


ALTER TABLE subscription_user_state ADD CONSTRAINT subscription_user_state_user_id_fk FOREIGN KEY (user_twitter_id) REFERENCES subscription_user (twitter_id) ON DELETE CASCADE;

ALTER TABLE subscription_user_state ADD CONSTRAINT subscription_user_state_unique UNIQUE (subscription_id, user_twitter_id);

COMMIT;
