CREATE OR REPLACE FUNCTION compute_next_task_pos ()
RETURNS pos_t
LANGUAGE 'sql'
AS $$
  SELECT (1 + COALESCE ((SELECT max ("position_in_list") FROM "task"), 0))::pos_t;
$$;

ALTER FUNCTION compute_next_task_pos ()
  OWNER TO "noda";
