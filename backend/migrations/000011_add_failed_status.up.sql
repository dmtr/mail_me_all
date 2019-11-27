ALTER TYPE subscription_status RENAME TO subscription_status_old;

CREATE TYPE subscription_status AS ENUM ('PREPARING', 'READY', 'SENDING', 'SENT', 'FAILED');

ALTER TABLE subscription_state ALTER COLUMN status TYPE subscription_status USING status::text::subscription_status; 

DROP TYPE subscription_status_old;
