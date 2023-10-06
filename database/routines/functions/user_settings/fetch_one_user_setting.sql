CREATE OR REPLACE FUNCTION fetch_one_user_setting (
  IN p_user_id     UUID,
  IN p_setting_key TEXT
)
RETURNS TABLE (
  "key"         "predefined_user_setting"."key"%TYPE,
  "description" "predefined_user_setting"."description"%TYPE,
  "value"       "predefined_user_setting"."default_value"%TYPE,
  "created_at"  "user_setting"."created_at"%TYPE,
  "updated_at"  "user_setting"."updated_at"%TYPE
)
LANGUAGE 'plpgsql'
AS $$
BEGIN
  CALL assert_user_exists (p_user_id);
  CALL assert_predefined_user_setting_exists (p_setting_key);
  RETURN QUERY
        SELECT us."key",
               df."description",
               us."value",
               us."created_at",
               us."updated_at"
          FROM "user_setting" us
    INNER JOIN "predefined_user_setting" df
            ON us."key" = df."key"
         WHERE us."user_id" = p_user_id AND
               us."key" = lower (p_setting_key);
END;
$$;

ALTER FUNCTION fetch_one_user_setting
      OWNER TO "noda";
