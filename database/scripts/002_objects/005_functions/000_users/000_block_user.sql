CREATE OR REPLACE FUNCTION block_user (
  IN p_user_id UUID
)
RETURNS BOOLEAN
AS $$
DECLARE
  is_already_blocked BOOLEAN;
  affected_rows INTEGER;
BEGIN
  CALL assert_user_exists (p_user_id);
  SELECT "is_blocked"
    INTO is_already_blocked
    FROM "user"
   WHERE "user_id" = p_user_id;

  IF is_already_blocked THEN
    RETURN FALSE;
  END IF;

  UPDATE "user"
     SET "is_blocked" = TRUE
   WHERE "user_id" = $1;

  GET DIAGNOSTICS affected_rows = ROW_COUNT;
  IF affected_rows > 0 THEN
    RETURN TRUE;
  END IF;
  RETURN FALSE;
END;
$$ LANGUAGE 'plpgsql';

ALTER FUNCTION block_user
      OWNER TO "noda";
