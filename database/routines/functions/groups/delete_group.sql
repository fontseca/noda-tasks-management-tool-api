CREATE OR REPLACE FUNCTION delete_group (
  IN p_owner_id UUID,
  IN p_group_id UUID
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_affected_rows INT;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_group_exists (p_group_id);
  DELETE FROM "group"
        WHERE "group_id" = p_group_id
          AND "owner_id" = p_owner_id;
  GET DIAGNOSTICS n_affected_rows := ROW_COUNT;
  IF n_affected_rows >= 1 THEN
    RETURN TRUE;
  END IF;
  RETURN FALSE;
END;
$$;

ALTER FUNCTION delete_group
      OWNER TO "noda";
