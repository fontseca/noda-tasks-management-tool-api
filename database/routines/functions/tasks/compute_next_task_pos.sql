CREATE OR REPLACE FUNCTION compute_next_task_pos
(
  IN p_list_scope "list"."list_id"%TYPE
)
RETURNS pos_t
LANGUAGE 'sql'
AS $$
  SELECT (1 + COALESCE ((SELECT max ("position_in_list")
                           FROM "task"
                          WHERE "list_id" = p_list_scope), 0))::pos_t;
$$;

ALTER FUNCTION compute_next_task_pos ("list"."list_id"%TYPE)
      OWNER TO "noda";
