CREATE OR REPLACE FUNCTION fetch_grouped_lists (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_group_id "list"."group_id"%TYPE,
  IN p_page      BIGINT,
  IN p_rpp       BIGINT,
  IN p_needle    TEXT,
  IN p_sort_expr TEXT
)
RETURNS SETOF "list"
LANGUAGE 'plpgsql'
AS $$
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_group_exists (p_owner_id, p_group_id);
  CALL validate_rpp_and_page (p_rpp, p_page);
  CALL gen_search_pattern (p_needle);
  CALL validate_sort_expr (p_sort_expr);
  RETURN QUERY
        SELECT *
          FROM "list" l
         WHERE l."is_archived" IS FALSE AND
               l."group_id" = p_group_id AND
               l."owner_id" = p_owner_id AND
               lower (concat (l."name", ' ', l."description")) ~ p_needle
      ORDER BY (CASE WHEN p_sort_expr = '' THEN l."created_at" END) DESC,
               (CASE WHEN p_sort_expr = '+name' THEN l."name" END) ASC,
               (CASE WHEN p_sort_expr = '-name' THEN l."name" END) DESC,
               (CASE WHEN p_sort_expr = '+description' THEN l."description" END) ASC,
               (CASE WHEN p_sort_expr = '-description' THEN l."description" END) DESC
         LIMIT p_rpp
        OFFSET (p_rpp * (p_page - 1));
END;
$$;

ALTER FUNCTION fetch_grouped_lists ("list"."owner_id"%TYPE,
                                    "list"."group_id"%TYPE,
                                    BIGINT,
                                    BIGINT,
                                    TEXT,
                                    TEXT)
      OWNER TO "noda";
