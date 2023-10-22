CREATE OR REPLACE PROCEDURE assert_group_exists (
  IN p_owner_id "group"."owner_id"%TYPE,
  IN p_group_id "group"."group_id"%TYPE
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_records INT;
  group_id_txt TEXT := p_group_id::TEXT;
BEGIN
  IF p_group_id IS NOT NULL THEN
    SELECT count(*)
      INTO n_records
      FROM "group"
     WHERE "owner_id" = p_owner_id AND
           "group_id" = p_group_id;
    IF n_records = 1 THEN
      RETURN;
    END IF;
  ELSE
    group_id_txt := '(NULL)';
  END IF;
  RAISE EXCEPTION 'nonexistent group with ID "%"', group_id_txt
       USING HINT = 'Please check the given group ID.';
END;
$$;

ALTER PROCEDURE assert_group_exists ("group"."owner_id"%TYPE,
                                     "group"."group_id"%TYPE)
       OWNER TO "noda";
