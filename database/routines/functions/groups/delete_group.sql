CREATE OR REPLACE FUNCTION delete_group (
  IN p_owner_id "group"."owner_id"%TYPE,
  IN p_group_id "group"."group_id"%TYPE
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_affected_rows INT;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_group_exists (p_owner_id, p_group_id);
  DELETE FROM "group"
        WHERE "group_id" = p_group_id
          AND "owner_id" = p_owner_id;
  GET DIAGNOSTICS n_affected_rows := ROW_COUNT;
  RETURN n_affected_rows;
END;
$$;

ALTER FUNCTION delete_group ("group"."owner_id"%TYPE,
                             "group"."group_id"%TYPE)
      OWNER TO "noda";
