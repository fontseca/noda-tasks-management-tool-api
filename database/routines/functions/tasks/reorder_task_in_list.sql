CREATE OR REPLACE FUNCTION reorder_task_in_list
(
  IN p_owner_id      "task"."owner_id"%TYPE,
  IN p_list_id       "task"."list_id"%TYPE,
  IN p_task_id       "task"."task_id"%TYPE,
  IN target_position pos_t
)
RETURNS BOOLEAN
RETURNS NULL ON NULL INPUT
LANGUAGE 'plpgsql'
AS $$
DECLARE
  affected_task_id "task"."task_id"%TYPE;
  obsolete_task_position pos_t := 1;
  affected_rows INTEGER;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  CALL assert_task_exists (p_owner_id, p_list_id, p_task_id);
  IF compute_next_task_pos (p_list_id) <= target_position THEN
    RETURN FALSE;
  END IF;
  SELECT t."task_id",
         t."position_in_list"
    INTO affected_task_id
    FROM "task" t
   WHERE t."owner_id" = p_owner_id
     AND t."list_id" = p_list_id
     AND t."position_in_list" = target_position;
  SELECT t."position_in_list"
    INTO obsolete_task_position
    FROM "task" t
   WHERE t."owner_id" = p_owner_id
     AND t."list_id" = p_list_id
     AND t."task_id" = p_task_id;
  IF target_position = obsolete_task_position THEN
    RETURN FALSE;
  END IF;
  /* Current task.  */
  UPDATE "task"
     SET "position_in_list" = target_position,
         "updated_at" = now ()
   WHERE "task"."owner_id" = p_owner_id
     AND "task"."list_id" = p_list_id
     AND "task"."task_id" = p_task_id;
  /* Affected task.  */
  UPDATE "task"
     SET "position_in_list" = obsolete_task_position,
         "updated_at" = now ()
   WHERE "task"."owner_id" = p_owner_id
     AND "task"."list_id" = p_list_id
     AND "task"."task_id" = affected_task_id;
  GET DIAGNOSTICS affected_rows := ROW_COUNT;
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION reorder_task_in_list ("task"."owner_id"%TYPE,
                                     "task"."list_id"%TYPE,
                                     "task"."task_id"%TYPE,
                                     pos_t)
OWNER TO "noda";
