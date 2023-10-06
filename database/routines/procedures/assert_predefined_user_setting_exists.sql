CREATE OR REPLACE PROCEDURE assert_predefined_user_setting_exists (
  IN p_setting_key TEXT
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_records INT;
  p_setting_key_txt TEXT := p_setting_key::TEXT;
BEGIN
  IF p_setting_key IS NOT NULL THEN
    SELECT count(*)
      INTO n_records
      FROM "predefined_user_setting"
    WHERE "key" = lower (p_setting_key);
    IF n_records = 1 THEN
      RETURN;
    END IF;
  ELSE
    p_setting_key_txt := '(NULL)';
  END IF;
  RAISE EXCEPTION 'nonexistent predefined user setting key: "%"', p_setting_key_txt
     USING DETAIL = 'No predefined user setting found with key "' || p_setting_key_txt || '".',
             HINT = 'Please check the given predefined user setting key.';
END;
$$;

ALTER PROCEDURE assert_predefined_user_setting_exists
       OWNER TO "noda";
