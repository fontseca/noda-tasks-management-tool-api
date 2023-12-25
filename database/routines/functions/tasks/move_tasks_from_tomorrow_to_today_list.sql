CREATE OR REPLACE FUNCTION move_tasks_from_tomorrow_to_today_list
(
  IN p_owner_id "task"."owner_id"%TYPE
)
RETURNS INTEGER
RETURNS NULL ON NULL INPUT
LANGUAGE 'plpgsql'
AS $$
DECLARE
  today_list_id "task"."list_id"%TYPE;
  tomorrow_list_id "task"."list_id"%TYPE;
  n_updated_tasks INTEGER;
BEGIN
  CALL assert_user_exists (p_owner_id);
  today_list_id := get_today_list_id (p_owner_id);
  tomorrow_list_id := get_tomorrow_list_id (p_owner_id);
  IF today_list_id IS NULL
    OR tomorrow_list_id IS NULL
  THEN
    RETURN 0;
  END IF;
  WITH "moved_tasks" AS
  (
       UPDATE "task"
          SET "list_id" = today_list_id
        WHERE "task"."owner_id" = p_owner_id
          AND "task"."list_id" = tomorrow_list_id
    RETURNING "task"."task_id"
  )
  SELECT count (*)
    INTO n_updated_tasks
    FROM "moved_tasks";
  RETURN n_updated_tasks;
END;
$$;

ALTER FUNCTION move_tasks_from_tomorrow_to_today_list ("task"."owner_id"%TYPE)
      OWNER TO "noda";
