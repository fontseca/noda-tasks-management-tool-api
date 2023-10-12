CREATE OR REPLACE PROCEDURE assert_group_exists (
  IN p_group_id UUID
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
    WHERE "group_id" = p_group_id;
    IF n_records = 1 THEN
      RETURN;
    END IF;
  ELSE
    group_id_txt := '(NULL)';
  END IF;
  RAISE EXCEPTION 'nonexistent group with ID "%"', group_id_txt
       USING HINT = 'Please check the given group ID';
END;
$$;

ALTER PROCEDURE assert_group_exists
       OWNER TO "noda";
