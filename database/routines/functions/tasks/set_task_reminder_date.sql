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
  task_due_date TIMESTAMPTZ;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  CALL assert_task_exists (p_owner_id, p_list_id, p_task_id);
  IF p_remind_at <= now () THEN
    RETURN FALSE;
  END IF;
  SELECT t."due_date"
    INTO task_due_date
    FROM "task" t
   WHERE t."owner_id" = p_owner_id
     AND t."list_id" = p_list_id
     AND t."task_id" = p_task_id;
  IF task_due_date IS NOT NULL AND p_remind_at >= task_due_date THEN
    RETURN FALSE;
  END IF;
  IF p_remind_at =
  (
    SELECT t."remind_at"
      FROM "task" t
     WHERE t."owner_id" = p_owner_id
       AND t."list_id" = p_list_id
       AND t."task_id" = p_task_id
  )
  THEN
    RETURN FALSE;
  END IF;
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
