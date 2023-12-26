CREATE OR REPLACE FUNCTION move_task_to_tomorrow_list
(
  IN p_owner_id "task"."owner_id"%TYPE,
  IN p_task_id  "task"."task_id"%TYPE
)
RETURNS BOOLEAN
RETURNS NULL ON NULL INPUT
LANGUAGE 'plpgsql'
AS $$
DECLARE
  tomorrow_list_id "task"."list_id"%TYPE;
  n_affected_rows INTEGER;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_task_exists_somewhere (p_owner_id, p_task_id);
  tomorrow_list_id := get_tomorrow_list_id (p_owner_id);
  IF tomorrow_list_id IS NULL THEN
    RETURN FALSE;
  END IF;
  UPDATE "task"
     SET "list_id" = tomorrow_list_id
   WHERE "task"."owner_id" = p_owner_id
     AND "task"."task_id" = p_task_id;
  GET DIAGNOSTICS n_affected_rows := ROW_COUNT;
  RETURN n_affected_rows;
END;
$$;

ALTER FUNCTION move_task_to_tomorrow_list ("task"."owner_id"%TYPE,
                                           "task"."task_id"%TYPE)
      OWNER TO "noda";
