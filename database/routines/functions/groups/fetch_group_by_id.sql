CREATE OR REPLACE FUNCTION fetch_group_by_id (
  IN p_owner_id UUID,
  IN p_group_id UUID
)
RETURNS SETOF "group"
LANGUAGE 'plpgsql'
AS $$
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_group_exists (p_group_id);
  RETURN QUERY
        SELECT *
          FROM "group"
         WHERE "group_id" = p_group_id AND
               "owner_id" = p_owner_id AND
               "is_archived" IS FALSE;
END;
$$;

ALTER FUNCTION fetch_group_by_id
      OWNER TO "noda";
