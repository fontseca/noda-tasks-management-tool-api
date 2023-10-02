CREATE OR REPLACE PROCEDURE assert_user_exists (
  IN p_user_id UUID
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_records INT;
BEGIN
  SELECT count(*)
    INTO n_records
    FROM "user"
   WHERE "user_id" = p_user_id;
  IF n_records <> 1 THEN
    RAISE EXCEPTION 'Nonexistent user ID.'
       USING DETAIL = 'No user found with the ID "' || p_user_id || '".',
               HINT = 'Please check the given user ID';
  END IF;
END;
$$;

ALTER PROCEDURE assert_user_exists
       OWNER TO "noda";
