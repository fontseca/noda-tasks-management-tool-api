CREATE OR REPLACE PROCEDURE assert_list_exists_somewhere (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_list_id  "list"."list_id"%TYPE
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_records INT;
  list_id_txt TEXT := p_list_id::TEXT;
BEGIN
  IF p_list_id IS NOT NULL THEN
    SELECT count (*)
      INTO n_records
      FROM "list"
     WHERE "owner_id" = p_owner_id AND
           "list_id" = p_list_id;
    IF n_records = 1 THEN
      RETURN;
    END IF;
  ELSE
    list_id_txt := '(NULL)';
  END IF;
  RAISE EXCEPTION 'nonexistent list with ID "%"', list_id_txt
       USING HINT = 'Please check the given list ID.';
END;
$$;

ALTER PROCEDURE assert_list_exists_somewhere ("list"."owner_id"%TYPE,
                                              "list"."list_id"%TYPE)
       OWNER TO "noda";
