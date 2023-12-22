CREATE OR REPLACE FUNCTION set_task_as_uncompleted
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
  affected_rows INT;
BEGIN
  IF (SELECT t."status"
        FROM "task" t
       WHERE t."owner_id" = p_owner_id
         AND t."list_id" = p_list_id
         AND t."task_id" = p_task_id) = 'unfinished'::task_status_t
  THEN
    RETURN FALSE;
  END IF;
  UPDATE "task"
     SET "status" = 'unfinished',
         "updated_at" = now ()
   WHERE "task"."owner_id" = p_owner_id
     AND "task"."list_id" = p_list_id
     AND "task"."task_id" = p_task_id;
  GET DIAGNOSTICS affected_rows := ROW_COUNT;
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION set_task_as_uncompleted ("task"."owner_id"%TYPE,
                                        "task"."task_id"%TYPE,
                                        "task"."task_id"%TYPE)
      OWNER TO "noda";
