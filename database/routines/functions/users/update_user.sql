CREATE OR REPLACE FUNCTION update_user (
  IN p_user_id UUID,
  IN p_first_name TEXT,
  IN p_middle_name TEXT,
  IN p_last_name TEXT,
  IN p_surname TEXT
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  rows_affected INT;
BEGIN
  CALL assert_user_exists (p_user_id);
  UPDATE "user" u
     SET u."first_name" = COALESCE (NULLIF (trim (p_first_name), ''), u."first_name"),
         u."middle_name" = COALESCE (NULLIF (trim (p_middle_name), ''), u."middle_name"),
         u."last_name" = COALESCE (NULLIF (trim (p_last_name), ''), u."last_name"),
         u."surname" = COALESCE (NULLIF (trim (p_surname), ''), u."surname"),
         u."updated_at" = 'now ()'
   WHERE u."user_id" = p_user_id;
  GET DIAGNOSTICS rows_affected = ROW_COUNT;
  RETURN rows_affected;
  IF rows_affected > 0 THEN
    RETURN TRUE;
  ELSE
    RETURN FALSE;
  END IF;
END;
$$;

ALTER FUNCTION update_user
      OWNER TO "noda";
