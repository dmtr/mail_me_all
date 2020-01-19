ALTER TYPE email_status RENAME TO email_status_old;

CREATE TYPE email_status AS ENUM ('NEW', 'SENT', 'CONFIRMED');

ALTER TABLE user_email_m2m ALTER COLUMN status TYPE email_status USING status::text::email_status; 

DROP TYPE email_status_old;
