CREATE OR REPLACE FUNCTION delete_list (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_group_id "list"."group_id"%TYPE,
  IN p_list_id  "list"."list_id"%TYPE
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_affected_rows INT;
  is_scattered_list CONSTANT BOOLEAN := p_group_id IS NULL;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_is_not_special_list (p_owner_id, p_list_id);
  CALL assert_list_exists (p_owner_id, p_group_id, p_list_id);
  DELETE FROM "list"
        WHERE "owner_id" = p_owner_id AND
              CASE WHEN is_scattered_list
                   THEN "group_id" IS NULL
                   ELSE "group_id" = p_group_id
               END AND
              "list_id" = p_list_id;
  GET DIAGNOSTICS n_affected_rows := ROW_COUNT;
  RETURN n_affected_rows;
END;
$$;

ALTER FUNCTION delete_list ("list"."owner_id"%TYPE,
                            "list"."group_id"%TYPE,
                            "list"."list_id"%TYPE)
      OWNER TO "noda";
