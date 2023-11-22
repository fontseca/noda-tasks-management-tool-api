CREATE OR REPLACE FUNCTION fetch_list_by_id (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_group_id "list"."group_id"%TYPE,
  IN p_list_id  "list"."list_id"%TYPE
)
RETURNS SETOF "list"
LANGUAGE 'plpgsql'
AS $$
DECLARE
  is_scattered_list CONSTANT BOOLEAN := p_group_id IS NULL;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_is_not_special_list (p_owner_id, p_list_id);
  CALL assert_list_exists (p_owner_id, p_group_id, p_list_id);
RETURN QUERY
      SELECT *
        FROM "list" l
       WHERE l."owner_id" = p_owner_id AND
             CASE WHEN is_scattered_list
                  THEN l."group_id" IS NULL
                  ELSE l."group_id" = "p_group_id"
             END AND
             l."list_id" = p_list_id;
END;
$$;

ALTER FUNCTION fetch_list_by_id ("list"."owner_id"%TYPE,
                                 "list"."group_id"%TYPE,
                                 "list"."list_id"%TYPE)
      OWNER TO "noda";
