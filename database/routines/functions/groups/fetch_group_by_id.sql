CREATE OR REPLACE FUNCTION fetch_group_by_id
(
  IN p_owner_id "group"."owner_id"%TYPE,
  IN p_group_id "group"."group_id"%TYPE
)
RETURNS SETOF "group"
LANGUAGE 'plpgsql'
AS $$
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_group_exists (p_owner_id, p_group_id);
  RETURN QUERY
        SELECT *
          FROM "group"
         WHERE "group_id" = p_group_id AND
               "owner_id" = p_owner_id;
END;
$$;

ALTER FUNCTION fetch_group_by_id ("group"."owner_id"%TYPE,
                                  "group"."group_id"%TYPE)
      OWNER TO "noda";
