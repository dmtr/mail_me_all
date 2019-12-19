BEGIN;

CREATE TABLE user_email_m2m(
    user_id UUID NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    CONSTRAINT user__email_m2m_user_account_id_fk FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE
);

COMMIT;
