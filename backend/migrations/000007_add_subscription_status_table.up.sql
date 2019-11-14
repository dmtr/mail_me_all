BEGIN;

CREATE TYPE subscription_status AS ENUM ('PREPARING', 'READY', 'SENDING', 'SENT');

CREATE TABLE subscription_state (
    id SERIAL PRIMARY KEY,
    status subscription_status NOT NULL,
    last_tweet_id VARCHAR,
    subscription_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT subscription_state_subscription_id_fk FOREIGN KEY (subscription_id) REFERENCES subscription (id) ON DELETE CASCADE
);

CREATE TRIGGER update_subscription_state
      before update
      on subscription_state
      for each row
      execute procedure update_timestamp()
  ;

CREATE TABLE tweet (
    id SERIAL PRIMARY KEY,
    tweet_id VARCHAR NOT NULL UNIQUE,
    tweet JSONB NOT NULL
);

CREATE TABLE subscription_state_tweet_m2m (
    subscription_state_id INTEGER NOT NULL,
    tweet_id INTEGER NOT NULL,
    CONSTRAINT subscription_state_tweet_m2m_subscription_state_id_fk FOREIGN KEY (subscription_state_id) REFERENCES subscription_state (id) ON DELETE CASCADE,
    CONSTRAINT subscription_state_tweet_m2m_tweet_id_fk FOREIGN KEY (tweet_id) REFERENCES tweet (id) ON DELETE CASCADE
);

COMMIT;
