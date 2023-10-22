CREATE OR REPLACE FUNCTION unblock_user (IN p_user_id "user"."user_id"%TYPE)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  is_blocked BOOLEAN;
  affected_rows INTEGER;
BEGIN
  CALL assert_user_exists (p_user_id);
  SELECT "user"."is_blocked"
    INTO is_blocked
    FROM "user"
   WHERE "user_id" = p_user_id;
  IF is_blocked IS FALSE THEN
    RETURN FALSE;
  END IF;
  UPDATE "user"
     SET "is_blocked" = FALSE
   WHERE "user_id" = p_user_id;
  GET DIAGNOSTICS affected_rows = ROW_COUNT;
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION unblock_user ("user"."user_id"%TYPE)
      OWNER TO "noda";
