BEGIN;

CREATE OR REPLACE FUNCTION get_day_of_week(d timestamp with time zone) RETURNS weekday
      language plpgsql
  AS $$
  DECLARE day weekday;
  BEGIN
	SELECT INTO day
	 CASE WHEN extract(isodow from d) = 1 THEN 'monday'
	      WHEN extract(isodow from d) = 2 THEN 'tuesday'
	      WHEN extract(isodow from d) = 3 THEN 'wensday'
	      WHEN extract(isodow from d) = 4 THEN 'thursday'
	      WHEN extract(isodow from d) = 5 THEN 'friday'
	      WHEN extract(isodow from d) = 6 THEN 'saturday'
	      WHEN extract(isodow from d) = 7 THEN 'sunday'
	 END;
	 RETURN day;
  END;
  $$
  ;

COMMIT;
