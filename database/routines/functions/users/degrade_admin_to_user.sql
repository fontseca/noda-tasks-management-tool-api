CREATE OR REPLACE FUNCTION degrade_admin_to_user (
  IN p_user_id UUID
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  affected_rows INTEGER;
  actual_role INTEGER;
BEGIN
  CALL assert_user_exists (p_user_id);
  SELECT "role_id"
    INTO actual_role
    FROM "user"
   WHERE "user_id" = p_user_id;

  IF actual_role = 2 THEN
    RETURN FALSE;
  END IF;

  UPDATE "user"
     SET "role_id" = 2,
         "updated_at" = 'now()'
   WHERE "user_id" = p_user_id;

  GET DIAGNOSTICS affected_rows = ROW_COUNT;
  IF affected_rows > 0 THEN
    RETURN TRUE;
  END IF;
  RETURN FALSE;
END;
$$;

ALTER FUNCTION degrade_admin_to_user
      OWNER TO "noda";
