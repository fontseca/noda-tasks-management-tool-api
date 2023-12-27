CREATE OR REPLACE FUNCTION restore_task_from_trash
(
  IN p_owner_id "task"."task_id"%TYPE,
  IN p_list_id  "task"."list_id"%TYPE,
  IN p_task_id  "task"."task_id"%TYPE
)
RETURNS BOOLEAN
RETURNS NULL ON NULL INPUT
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_inserted_rows INTEGER;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  IF 0 = (SELECT count (1)
            FROM "trashed_task" t
           WHERE t."owner_id" = p_owner_id
             AND t."list_id" = p_list_id
             AND t."task_id" = p_task_id)
  THEN
    RETURN FALSE;
  END IF;
  WITH "moved_task" AS
  (
    DELETE FROM "trashed_task"
          WHERE "trashed_task"."owner_id" = p_owner_id
            AND "trashed_task"."list_id" = p_list_id
            AND "trashed_task"."task_id" = p_task_id
      RETURNING *
  )
  INSERT INTO "task" ("task_id",
                      "owner_id",
                      "list_id",
                      "position_in_list",
                      "title",
                      "headline",
                      "description",
                      "priority",
                      "status",
                      "is_pinned",
                      "due_date",
                      "remind_at",
                      "completed_at",
                      "created_at",
                      "updated_at")
       SELECT "task_id",
              "owner_id",
              "list_id",
              "position_in_list",
              "title",
              "headline",
              "description",
              "priority",
              "status",
              "is_pinned",
              "due_date",
              "remind_at",
              "completed_at",
              "created_at",
              "updated_at"
         FROM "moved_task";
  SELECT count (1)
    INTO n_inserted_rows
    FROM "task" t
   WHERE t."owner_id" = p_owner_id
     AND t."list_id" = p_list_id
     AND t."task_id" = p_task_id;
  RETURN n_inserted_rows = 1;
END;
$$;

ALTER FUNCTION restore_task_from_trash ("task"."task_id"%TYPE,
                                        "task"."list_id"%TYPE,
                                        "task"."task_id"%TYPE)
  OWNER TO "noda";
