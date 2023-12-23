CREATE OR REPLACE FUNCTION pin_task
(
  IN p_owner_id "task"."owner_id"%TYPE,
  IN p_list_id  "task"."task_id"%TYPE,
  IN p_task_id  "task"."task_id"%TYPE
)
RETURNS BOOLEAN
RETURNS NULL ON NULL INPUT
LANGUAGE 'plpgsql'
AS $$
DECLARE
  affected_rows INTEGER;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  CALL assert_task_exists (p_owner_id, p_list_id, p_task_id);
  IF (SELECT t."is_pinned"
        FROM "task" t
       WHERE t."owner_id" = p_owner_id
         AND t."list_id" = p_list_id
         AND t."task_id" = p_task_id)
  THEN
    RETURN FALSE;
  END IF;
  UPDATE "task"
     SET "is_pinned" = TRUE,
         "updated_at" = now ()
   WHERE "task"."owner_id" = p_owner_id
     AND "task"."list_id" = p_list_id
     AND "task"."task_id" = p_task_id;
  GET DIAGNOSTICS affected_rows := ROW_COUNT;
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION pin_task ("task"."owner_id"%TYPE,
                         "task"."task_id"%TYPE,
                         "task"."task_id"%TYPE)
      OWNER TO "noda";
