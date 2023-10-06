CREATE OR REPLACE FUNCTION delete_user_hardly (
  IN p_user_id UUID
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_affected_rows INT;
BEGIN
  CALL assert_user_exists (p_user_id);
  DELETE FROM "user"
        WHERE "user_id" = p_user_id;
  GET DIAGNOSTICS n_affected_rows := ROW_COUNT;
  IF n_affected_rows >= 1 THEN
    RETURN TRUE;
  END IF;
  RETURN FALSE;
END;
$$;

ALTER FUNCTION delete_user_hardly
      OWNER TO "noda";
