CREATE OR REPLACE FUNCTION fetch_task_by_id
(
  IN p_owner_id "task"."owner_id"%TYPE,
  IN p_list_id  "task"."list_id"%TYPE,
  IN p_task_id  "task"."task_id"%TYPE
)
RETURNS SETOF "task"
LANGUAGE 'plpgsql'
AS $$
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  CALL assert_task_exists (p_owner_id, p_list_id, p_task_id);
  RETURN QUERY
        SELECT *
          FROM "task" t
         WHERE t."owner_id" = p_owner_id AND
               t."list_id" = p_list_id AND
               t."task_id" = p_task_id;
END;
$$;

ALTER FUNCTION fetch_task_by_id ("task"."owner_id"%TYPE,
                                 "task"."list_id"%TYPE,
                                 "task"."task_id"%TYPE)
      OWNER TO "noda";
