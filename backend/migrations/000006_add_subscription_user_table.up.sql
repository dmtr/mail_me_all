BEGIN;

CREATE TABLE subscription_user(
    id SERIAL PRIMARY KEY,
    twitter_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    profile_image_url VARCHAR NOT NULL,
    screen_name VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_subscription_user
      before update
      on subscription_user
      for each row
      execute procedure update_timestamp()
  ;

CREATE TABLE subscription_user_m2m(
    subscription_id UUID NOT NULL,
    user_id INTEGER NOT NULL,
    CONSTRAINT subscription_user_m2m_subscription_id_fk FOREIGN KEY (subscription_id) REFERENCES subscription (id) ON DELETE CASCADE,
    CONSTRAINT subscription_user_m2m_user_id_fk FOREIGN KEY (user_id) REFERENCES subscription_user (id) ON DELETE CASCADE
);

COMMIT;
