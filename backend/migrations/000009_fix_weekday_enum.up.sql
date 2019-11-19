ALTER TYPE weekday RENAME TO weekday_old;

CREATE TYPE weekday AS ENUM ('monday', 'tuesday', 'wensday', 'thursday', 'friday', 'saturday', 'sunday');

ALTER TABLE subscription ALTER COLUMN day TYPE weekday USING day::text::weekday; 

DROP TYPE weekday_old;
