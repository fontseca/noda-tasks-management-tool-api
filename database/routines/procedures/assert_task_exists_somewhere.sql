CREATE OR REPLACE PROCEDURE assert_task_exists_somewhere
(
  IN p_owner_id "task"."owner_id"%TYPE,
  IN p_task_id  "task"."task_id"%TYPE
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_records INT;
  task_id_txt TEXT := p_task_id::TEXT;
BEGIN
  IF p_task_id IS NOT NULL THEN
    SELECT count (*)
      INTO n_records
      FROM "task"
     WHERE "owner_id" = p_owner_id AND
           "task_id" = p_task_id;
    IF n_records = 1 THEN
      RETURN;
    END IF;
  ELSE
    task_id_txt := '(NULL)';
  END IF;
  RAISE EXCEPTION 'nonexistent task with ID "%"', task_id_txt
       USING HINT = 'Please check the given task ID.';
END;
$$;

ALTER PROCEDURE assert_task_exists_somewhere ("task"."owner_id"%TYPE,
                                              "task"."task_id"%TYPE)
       OWNER TO "noda";
