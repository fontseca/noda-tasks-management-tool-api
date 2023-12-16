CREATE OR REPLACE FUNCTION fetch_tasks
(
  IN p_owner_id  "task"."owner_id"%TYPE,
  IN p_list_id   "task"."list_id"%TYPE,
  IN p_page      BIGINT,
  IN p_rpp       BIGINT,
  IN p_needle    TEXT,
  IN p_sort_expr TEXT
)
RETURNS SETOF "task"
RETURNS NULL ON NULL INPUT
LANGUAGE 'plpgsql'
AS $$
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  CALL validate_rpp_and_page (p_rpp, p_page);
  CALL gen_search_pattern (p_needle);
  CALL validate_sort_expr (p_sort_expr);
  RETURN QUERY
        SELECT *
          FROM "task" t
         WHERE t."owner_id" = p_owner_id
           AND t."list_id" = p_list_id
           AND lower (concat (t."title", '',
                              t."headline", '',
                              t."description")) ~ p_needle
      ORDER BY (CASE WHEN p_sort_expr = '' THEN t."position_in_list" END) ASC,
               (CASE WHEN p_sort_expr = '+position_in_list' THEN t."position_in_list" END) ASC,
               (CASE WHEN p_sort_expr = '-position_in_list' THEN t."position_in_list" END) DESC,
               (CASE WHEN p_sort_expr = '+title' THEN t."title" END) ASC,
               (CASE WHEN p_sort_expr = '-title' THEN t."title" END) DESC,
               (CASE WHEN p_sort_expr = '+headline' THEN t."headline" END) ASC,
               (CASE WHEN p_sort_expr = '-headline' THEN t."headline" END) DESC,
               (CASE WHEN p_sort_expr = '+description' THEN t."description" END) ASC,
               (CASE WHEN p_sort_expr = '-description' THEN t."description" END) DESC,
               (CASE WHEN p_sort_expr = '+priority' THEN t."priority" END) ASC,
               (CASE WHEN p_sort_expr = '-priority' THEN t."priority" END) DESC,
               (CASE WHEN p_sort_expr = '+status' THEN t."status" END) ASC,
               (CASE WHEN p_sort_expr = '-status' THEN t."status" END) DESC,
               (CASE WHEN p_sort_expr = '+is_pinned' THEN t."is_pinned" END) ASC,
               (CASE WHEN p_sort_expr = '-is_pinned' THEN t."is_pinned" END) DESC,
               (CASE WHEN p_sort_expr = '+due_date' THEN t."due_date" END) ASC,
               (CASE WHEN p_sort_expr = '-due_date' THEN t."due_date" END) DESC,
               (CASE WHEN p_sort_expr = '+remind_at' THEN t."remind_at" END) ASC,
               (CASE WHEN p_sort_expr = '-remind_at' THEN t."remind_at" END) DESC,
               (CASE WHEN p_sort_expr = '+completed_at' THEN t."completed_at" END) ASC,
               (CASE WHEN p_sort_expr = '-completed_at' THEN t."completed_at" END) DESC,
               (CASE WHEN p_sort_expr = '+created_at' THEN t."created_at" END) ASC,
               (CASE WHEN p_sort_expr = '-created_at' THEN t."created_at" END) DESC,
               (CASE WHEN p_sort_expr = '+updated_at' THEN t."updated_at" END) ASC,
               (CASE WHEN p_sort_expr = '-updated_at' THEN t."updated_at" END) DESC
         LIMIT p_rpp
        OFFSET (p_rpp * (p_page - 1));
END;
$$;

ALTER FUNCTION fetch_tasks ("task"."owner_id"%TYPE,
                            "task"."list_id"%TYPE,
                            BIGINT,
                            BIGINT,
                            TEXT,
                            TEXT)
      OWNER TO "noda";
