BEGIN;

ALTER TABLE user_account ADD COLUMN email VARCHAR;

ALTER TABLE user_account DROP COLUMN fb_token;

CREATE TABLE token (
    user_id UUID,
    fb_token VARCHAR NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE FUNCTION update_timestamp() RETURNS TRIGGER
      language plpgsql
  as $$
  BEGIN
	  NEW.updated_at = now();
	  RETURN NEW;
  END;
  $$
  ;

CREATE TRIGGER update_user_account
      before update
      on user_account
      for each row
      execute procedure update_timestamp()
  ;

CREATE TRIGGER update_token
      before update
      on token
      for each row
      execute procedure update_timestamp()
  ;

COMMIT;
