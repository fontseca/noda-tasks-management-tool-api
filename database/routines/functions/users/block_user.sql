CREATE OR REPLACE FUNCTION block_user (IN p_user_id "user"."user_id"%TYPE)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
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
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION block_user ("user"."user_id"%TYPE)
      OWNER TO "noda";
