CREATE OR REPLACE FUNCTION convert_to_scattered_list (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_list_id  "list"."list_id"%TYPE
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  current_group_id "list"."group_id"%TYPE;
  n_affected_rows INT;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_is_not_special_list (p_owner_id, p_list_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  SELECT l."group_id"
    INTO current_group_id
    FROM "list" l
   WHERE l."owner_id" = p_owner_id AND
         l."list_id" = p_list_id;
  IF current_group_id IS NULL THEN
    RETURN FALSE;
  END IF;
  UPDATE "list"
     SET "group_id" = NULL,
         "updated_at" = 'now ()'
   WHERE "owner_id" = p_owner_id AND
         "list_id" = p_list_id;
  GET DIAGNOSTICS n_affected_rows := ROW_COUNT;
  RETURN n_affected_rows;
END;
$$;

ALTER FUNCTION convert_to_scattered_list ("list"."owner_id"%TYPE,
                                          "list"."list_id"%TYPE)
      OWNER TO "noda";
