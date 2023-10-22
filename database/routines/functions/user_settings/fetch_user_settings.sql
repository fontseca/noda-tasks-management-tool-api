CREATE OR REPLACE FUNCTION fetch_user_settings (
  IN p_user_id   "user"."user_id"%TYPE,
  IN p_page      BIGINT,
  IN p_rpp       BIGINT,
  IN p_needle    TEXT,
  IN p_sort_expr TEXT
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
  CALL validate_rpp_and_page (p_rpp, p_page);
  CALL validate_sort_expr (p_sort_expr);
  CALL gen_search_pattern (p_needle);
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
               lower (concat (us."value", ' ', df."description")) ~ p_needle
      ORDER BY (CASE WHEN p_sort_expr = '' THEN us."created_at" END) DESC,
               (CASE WHEN p_sort_expr = '+key' THEN us."key" END) ASC,
               (CASE WHEN p_sort_expr = '-key' THEN us."key" END) DESC,
               (CASE WHEN p_sort_expr = '+description' THEN df."description" END) ASC,
               (CASE WHEN p_sort_expr = '-description' THEN df."description" END) DESC,
               (CASE WHEN p_sort_expr = '+created_at' THEN us."created_at" END) ASC,
               (CASE WHEN p_sort_expr = '-created_at' THEN us."created_at" END) DESC,
               (CASE WHEN p_sort_expr = '+updated_at' THEN us."updated_at" END) ASC,
               (CASE WHEN p_sort_expr = '-updated_at' THEN us."updated_at" END) DESC
         LIMIT p_rpp
        OFFSET (p_rpp * (p_page - 1));
END;
$$;

ALTER FUNCTION fetch_user_settings ("user"."user_id"%TYPE, BIGINT, BIGINT, TEXT, TEXT)
      OWNER TO "noda";
