CREATE OR REPLACE FUNCTION update_user_setting (
  IN p_user_id          "user"."user_id"%TYPE,
  IN p_user_setting_key "user_setting"."key"%TYPE,
  IN p_user_setting_val "user_setting"."value"%TYPE
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
         "updated_at" = 'now ()'
   WHERE "user_id" = p_user_id AND
         "key" = p_user_setting_key;
  GET DIAGNOSTICS affected_rows := ROW_COUNT;
  RAISE NOTICE 'count is %', affected_rows;
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION update_user_setting ("user"."user_id"%TYPE,
                                    "user_setting"."key"%TYPE,
                                    "user_setting"."value"%TYPE)
      OWNER TO "noda";
