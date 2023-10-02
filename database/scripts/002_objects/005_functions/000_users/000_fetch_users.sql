CREATE OR REPLACE FUNCTION fetch_users (
  IN p_rpp BIGINT, /* records per page  */
  IN p_page BIGINT,
  IN p_sort_expr TEXT,
  IN p_needle TEXT
)
RETURNS TABLE(
  "id"          "user"."user_id"%TYPE,
  "role"        "user"."role_id"%TYPE,
  "first_name"  "user"."first_name"%TYPE,
  "middle_name" "user"."middle_name"%TYPE,
  "last_name"   "user"."last_name"%TYPE,
  "surname"     "user"."surname"%TYPE,
  "picture_url" "user"."picture_url"%TYPE,
  "email"       "user"."email"%TYPE,
  "is_blocked"  "user"."is_blocked"%TYPE,
  "created_at"  "user"."created_at"%TYPE,
  "updated_at"  "user"."updated_at"%TYPE
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  max_int_64 CONSTANT BIGINT := 9223372036854775807;
  max_valid_page_before_overflow BIGINT;
  sort_expr TEXT := lower(p_sort_expr);
  needle_pattern CONSTANT TEXT := make_search_pattern (p_needle);
BEGIN
  IF p_sort_expr = '' OR p_sort_expr IS NULL THEN
    sort_expr := '';
  ELSIF left(sort_expr, 1) NOT IN ('+', '-') THEN
    sort_expr := '+' || sort_expr;
  END IF;

  /* Make sure we can retrieve `p_rpp' records in just
     `p_page' pages.  If not, then use the maximum value
     for a page.  */

  IF p_rpp <= 0 OR p_rpp IS NULL THEN
    p_rpp := 1;
  END IF;
  IF p_page <= 0 OR p_page IS NULL THEN
    p_page := 1;
  END IF;
  max_valid_page_before_overflow := (max_int_64 / p_rpp) - 1;
  IF p_page > max_valid_page_before_overflow THEN
    p_page := max_valid_page_before_overflow;
  END IF;

  RETURN QUERY
        SELECT u."user_id" AS "id",
               u."role_id" AS "role",
               u."first_name",
               u."middle_name",
               u."last_name",
               u."surname",
               u."picture_url",
               u."email",
               u."is_blocked",
               u."created_at",
               u."updated_at"
          FROM "user" u
         WHERE u."is_blocked" IS FALSE
               AND lower(concat(
                 u."first_name", ' ',
                 u."middle_name", ' ',
                 u."last_name", ' ',
                 u."surname")) ~ needle_pattern
      ORDER BY (CASE WHEN sort_expr = '' THEN u."created_at" END) DESC,
               (CASE WHEN sort_expr = '+user_id' THEN u."user_id" END) ASC,
               (CASE WHEN sort_expr = '-user_id' THEN u."user_id" END) DESC,
               (CASE WHEN sort_expr = '+first_name' THEN u."first_name" END) ASC,
               (CASE WHEN sort_expr = '-first_name' THEN u."first_name" END) DESC,
               (CASE WHEN sort_expr = '+middle_name' THEN u."middle_name" END) ASC,
               (CASE WHEN sort_expr = '-middle_name' THEN u."middle_name" END) DESC,
               (CASE WHEN sort_expr = '+last_name' THEN u."last_name" END) ASC,
               (CASE WHEN sort_expr = '-last_name' THEN u."last_name" END) DESC,
               (CASE WHEN sort_expr = '+surname' THEN u."surname" END) ASC,
               (CASE WHEN sort_expr = '-surname' THEN u."surname" END) DESC,
               (CASE WHEN sort_expr = '+email' THEN u."email" END) ASC,
               (CASE WHEN sort_expr = '-email' THEN u."email" END) DESC,
               (CASE WHEN sort_expr = '+created_at' THEN u."created_at" END) ASC,
               (CASE WHEN sort_expr = '-created_at' THEN u."created_at" END) DESC,
               (CASE WHEN sort_expr = '+updated_at' THEN u."updated_at" END) ASC,
               (CASE WHEN sort_expr = '-updated_at' THEN u."updated_at" END) DESC
         LIMIT p_rpp
        OFFSET (p_rpp * (p_page - 1));
END;
$$;

ALTER FUNCTION fetch_users
      OWNER TO "noda";
