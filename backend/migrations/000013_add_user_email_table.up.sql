BEGIN;

CREATE TYPE email_status AS ENUM ('NEW', 'CONFIRMED');

CREATE TABLE user_email_m2m(
    user_id UUID NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    status email_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_email_m2m_user_account_id_fk FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE TRIGGER update_user_email_m2m
      before update
      on user_email_m2m
      for each row
      execute procedure update_timestamp()
  ;

COMMIT;
