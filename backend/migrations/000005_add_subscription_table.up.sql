BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE weekday AS ENUM ('monday', 'tuedasy', 'wensday', 'thursday', 'friday', 'saturday', 'sunday');

CREATE TABLE subscription(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    title VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    day weekday,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT subscription_user_id_fk FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE TRIGGER update_subscription
      before update
      on subscription
      for each row
      execute procedure update_timestamp()
  ;

COMMIT;
