CREATE OR REPLACE FUNCTION fetch_lists (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_page      BIGINT,
  IN p_rpp       BIGINT,
  IN p_needle    TEXT,
  IN p_sort_expr TEXT
)
RETURNS SETOF "list"
LANGUAGE 'plpgsql'
AS $$
DECLARE
  today_list_id "list"."list_id"%TYPE;
  tomorrow_list_id "list"."list_id"%TYPE;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL validate_rpp_and_page (p_rpp, p_page);
  CALL gen_search_pattern (p_needle);
  CALL validate_sort_expr (p_sort_expr);
  today_list_id := get_today_list_id (p_owner_id);
  tomorrow_list_id := get_tomorrow_list_id (p_owner_id);
  RETURN QUERY
        SELECT *
          FROM "list" l
         WHERE l."is_archived" IS FALSE AND
               l."owner_id" = p_owner_id AND
               CASE WHEN today_list_id IS NULL
                    THEN TRUE
                    ELSE l."list_id" <> today_list_id
                END AND
               CASE WHEN tomorrow_list_id IS NULL
                    THEN TRUE
                    ELSE l."list_id" <> tomorrow_list_id
                END AND
               lower (concat (l."name", ' ', l."description")) ~ p_needle
      ORDER BY (CASE WHEN p_sort_expr = '' THEN l."created_at" END) DESC,
               (CASE WHEN p_sort_expr = '+name' THEN l."name" END) ASC,
               (CASE WHEN p_sort_expr = '-name' THEN l."name" END) DESC,
               (CASE WHEN p_sort_expr = '+description' THEN l."description" END) ASC,
               (CASE WHEN p_sort_expr = '-description' THEN l."description" END) DESC,
               (CASE WHEN p_sort_expr = '+is_archived' THEN l."is_archived" END) ASC,
               (CASE WHEN p_sort_expr = '-is_archived' THEN l."is_archived" END) DESC,
               (CASE WHEN p_sort_expr = '+archived_at' THEN l."archived_at" END) ASC,
               (CASE WHEN p_sort_expr = '-archived_at' THEN l."archived_at" END) DESC,
               (CASE WHEN p_sort_expr = '+created_at' THEN l."created_at" END) ASC,
               (CASE WHEN p_sort_expr = '-created_at' THEN l."created_at" END) DESC,
               (CASE WHEN p_sort_expr = '+updated_at' THEN l."updated_at" END) ASC,
               (CASE WHEN p_sort_expr = '-updated_at' THEN l."updated_at" END) DESC
         LIMIT p_rpp
        OFFSET (p_rpp * (p_page - 1));
END;
$$;

ALTER FUNCTION fetch_lists ("list"."owner_id"%TYPE,
                            BIGINT,
                            BIGINT,
                            TEXT,
                            TEXT)
      OWNER TO "noda";
