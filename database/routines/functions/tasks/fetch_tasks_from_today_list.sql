CREATE OR REPLACE FUNCTION fetch_tasks_from_today_list
(
  IN p_owner_id  "task"."owner_id"%TYPE,
  IN p_page      BIGINT,
  IN p_rpp       BIGINT,
  IN p_needle    TEXT,
  IN p_sort_expr TEXT
)
RETURNS SETOF "task"
RETURNS NULL ON NULL INPUT
LANGUAGE 'plpgsql'
AS $$
DECLARE
  today_list_id "list"."list_id"%TYPE;
BEGIN
  SELECT get_today_list_id (p_owner_id)
    INTO today_list_id;
  IF today_list_id IS NULL THEN
    RETURN;
  END IF;
  RETURN QUERY
        SELECT *
          FROM fetch_tasks (p_owner_id,
                            today_list_id,
                            p_page,
                            p_rpp,
                            p_needle,
                            p_sort_expr);
END;
$$;

ALTER FUNCTION fetch_tasks_from_today_list ("task"."owner_id"%TYPE,
                                            BIGINT,
                                            BIGINT,
                                            TEXT,
                                            TEXT)
      OWNER TO "noda";
