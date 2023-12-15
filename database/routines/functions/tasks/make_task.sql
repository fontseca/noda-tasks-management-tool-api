CREATE OR REPLACE FUNCTION make_task
(
  IN p_owner_id      "user"."user_id"%TYPE,
  IN p_list_id       "task"."list_id"%TYPE,
  IN p_task_creation task_creation_t
)
RETURNS "task"."task_id"%TYPE
LANGUAGE 'plpgsql'
AS $$
DECLARE
  inserted_id "task"."task_id"%TYPE;
  actual_list_id "task"."task_id"%TYPE := p_list_id;
  n_similar_titles INT := 0;
  actual_list_title "task"."title"%TYPE := p_task_creation."title";
BEGIN
  CALL assert_user_exists (p_owner_id);
  IF actual_list_id IS NOT NULL THEN
    CALL assert_list_exists (p_owner_id, NULL, actual_list_id);
  ELSE
    actual_list_id := get_today_list_id (p_owner_id);
    IF actual_list_id IS NULL THEN
      actual_list_id := make_today_list (p_owner_id);
    END IF;
  END IF;
  SELECT count (*)
    INTO n_similar_titles
    FROM "task" t
   WHERE "list_id" = p_list_id AND
         regexp_count (t."title", '^' || quote_meta (actual_list_title) || '(?: \(\d+\))?$') = 1;
  IF n_similar_titles > 0 THEN
    actual_list_title := concat (actual_list_title, ' ' , '(', n_similar_titles, ')');
  END IF;
  INSERT INTO "task" ("owner_id",
                      "list_id",
                      "position_in_list",
                      "title",
                      "headline",
                      "description",
                      "priority",
                      "status",
                      "due_date",
                      "remind_at")
       VALUES (p_owner_id,
               actual_list_id,
               compute_next_task_pos (p_list_id),
               actual_list_title,
               NULLIF (p_task_creation."headline", ''),
               NULLIF (p_task_creation."description", ''),
               COALESCE (p_task_creation."priority", 'normal'::task_priority_t),
               COALESCE (p_task_creation."status", 'in progress'::task_status_t),
               p_task_creation."due_date",
               p_task_creation."remind_at")
    RETURNING "task_id"
         INTO inserted_id;
  RETURN inserted_id;
END;
$$;

ALTER FUNCTION make_task ("user"."user_id"%TYPE,
                          "task"."list_id"%TYPE,
                          task_creation_t)
      OWNER TO "noda";
