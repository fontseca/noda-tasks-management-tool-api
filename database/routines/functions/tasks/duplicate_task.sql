CREATE OR REPLACE FUNCTION duplicate_task
(
  IN p_owner_id "task"."owner_id"%TYPE,
  IN p_task_id  "task"."task_id"%TYPE
)
RETURNS "list"."list_id"%TYPE
LANGUAGE 'plpgsql'
AS $$
DECLARE
  current_task "task"%ROWTYPE;
  new_task_id "task"."task_id"%TYPE;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_task_exists_somewhere (p_owner_id, p_task_id);
  SELECT *
    INTO current_task
    FROM "task"
   WHERE "owner_id" = p_owner_id AND
         "task_id" = p_task_id;
  new_task_id := make_task (p_owner_id,
                            current_task."list_id",
                            ROW(current_task."title",
                                current_task."headline",
                                current_task."description",
                                current_task."priority",
                                current_task."status",
                                current_task."due_date",
                                current_task."remind_at"));
  /* TODO: Duplicate all steps and attachments.  */
  RETURN new_task_id;
END;
$$;

ALTER FUNCTION duplicate_task("task"."owner_id"%TYPE,
                              "task"."task_id"%TYPE)
      OWNER TO "noda";
