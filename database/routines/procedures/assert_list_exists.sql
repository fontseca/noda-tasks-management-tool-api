CREATE OR REPLACE PROCEDURE assert_list_exists (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_group_id "list"."group_id"%TYPE,
  IN p_list_id  "list"."list_id"%TYPE
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_records INT;
  list_id_txt TEXT := p_list_id::TEXT;
  is_grouped_list CONSTANT BOOLEAN := p_group_id IS NOT NULL;
  hint TEXT := 'Please check the given list ID.';
BEGIN
  IF is_grouped_list THEN
    CALL assert_group_exists(p_owner_id, p_group_id);
  END IF;
  IF p_list_id IS NOT NULL THEN
    SELECT count (*)
      INTO n_records
      FROM "list"
     WHERE "owner_id" = p_owner_id AND
           CASE WHEN is_grouped_list
                THEN "group_id" = p_group_id
                ELSE "group_id" IS NULL
           END AND
           "list_id" = p_list_id;
    IF n_records = 1 THEN
      RETURN;
    END IF;
  ELSE
    list_id_txt := '(NULL)';
  END IF;
  IF NOT is_grouped_list THEN
    hint := 'Please check the given list ID or whether the list is scattered.';
  END IF;
  RAISE EXCEPTION 'nonexistent list with ID "%"', list_id_txt
       USING HINT = hint;
END;
$$;

ALTER PROCEDURE assert_list_exists ("list"."owner_id"%TYPE,
                                    "list"."group_id"%TYPE,
                                    "list"."list_id"%TYPE)
       OWNER TO "noda";
