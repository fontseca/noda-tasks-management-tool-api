CREATE OR REPLACE FUNCTION set_task_reminder_date
(
  IN p_owner_id  "task"."owner_id"%TYPE,
  IN p_list_id   "task"."task_id"%TYPE,
  IN p_task_id   "task"."task_id"%TYPE,
  IN p_remind_at TIMESTAMPTZ
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  affected_rows INTEGER;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  CALL assert_task_exists (p_owner_id, p_list_id, p_task_id);
  UPDATE "task"
     SET "remind_at" = p_remind_at,
         "updated_at" = now ()
  WHERE "task"."owner_id" = p_owner_id
    AND "task"."list_id" = p_list_id
    AND "task"."task_id" = p_task_id;
  GET DIAGNOSTICS affected_rows := ROW_COUNT;
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION set_task_reminder_date ("task"."owner_id"%TYPE,
                                       "task"."task_id"%TYPE,
                                       "task"."task_id"%TYPE,
                                       TIMESTAMPTZ)
      OWNER TO "noda";
