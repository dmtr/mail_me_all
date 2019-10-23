BEGIN;

ALTER TABLE user_account DROP COLUMN fb_id;

DROP TABLE token;

CREATE TABLE social_account (
    user_id UUID,
    social_account_id VARCHAR NOT NULL,
    access_token VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE TABLE tw_account (
    user_id UUID,
    token_secret VARCHAR NOT NULL,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE
) INHERITS(social_account) ;

CREATE TRIGGER update_tw_account
      before update
      on tw_account
      for each row
      execute procedure update_timestamp()
  ;

COMMIT;
