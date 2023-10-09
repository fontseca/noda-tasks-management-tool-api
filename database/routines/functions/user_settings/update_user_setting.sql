CREATE OR REPLACE FUNCTION update_user_setting (
  IN p_user_id          UUID,
  IN p_user_setting_key TEXT,
  IN p_user_setting_val JSON
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  actual_setting_val TEXT;
  affected_rows INT;
BEGIN
  CALL assert_user_exists (p_user_id);
  CALL assert_predefined_user_setting_exists (p_user_setting_key);

  SELECT "value"
    INTO actual_setting_val
    FROM "user_setting"
   WHERE "user_id" = p_user_id AND
         "key" = p_user_setting_key;

  IF p_user_setting_val::TEXT = actual_setting_val::TEXT THEN
    RETURN FALSE;
  END IF;

  UPDATE "user_setting"
     SET "value" = p_user_setting_val,
         "updated_at" = 'now()'
   WHERE "user_id" = p_user_id AND
         "key" = p_user_setting_key;
  GET DIAGNOSTICS affected_rows := ROW_COUNT;
  RAISE NOTICE 'count is %', affected_rows;
  IF affected_rows >= 1 THEN
    RETURN TRUE;
  END IF;
  RETURN FALSE;
END;
$$;

ALTER FUNCTION update_user_setting
      OWNER TO "noda";
