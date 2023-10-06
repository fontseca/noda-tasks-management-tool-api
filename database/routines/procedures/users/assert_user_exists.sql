CREATE OR REPLACE PROCEDURE assert_user_exists (
  IN p_user_id UUID
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_records INT;
  user_id_txt TEXT := p_user_id::TEXT;
BEGIN
  IF p_user_id IS NOT NULL THEN
    SELECT count(*)
      INTO n_records
      FROM "user"
    WHERE "user_id" = p_user_id;
    IF n_records = 1 THEN
      RETURN;
    END IF;
  ELSE
    user_id_txt := '(NULL)';
  END IF;
  RAISE EXCEPTION 'nonexistent user with ID "%"', user_id_txt
       USING HINT = 'Please check the given user ID';
END;
$$;

ALTER PROCEDURE assert_user_exists
       OWNER TO "noda";
