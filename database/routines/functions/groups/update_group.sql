CREATE OR REPLACE FUNCTION update_group (
  IN p_owner_id "group"."owner_id"%TYPE,
  IN p_group_id "group"."group_id"%TYPE,
  IN p_g_name   "group"."name"%TYPE,
  IN p_g_desc   "group"."description"%TYPE
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_affected_rows INT;
  old_g_name "group"."name"%TYPE;
  old_g_desc "group"."description"%TYPE;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_group_exists (p_group_id);
  SELECT g."name",
         g."description"
    INTO old_g_name,
         old_g_desc
    FROM "group" g
   WHERE g."group_id" = p_group_id AND
         g."owner_id" = p_owner_id;
  IF (old_g_name = p_g_name OR p_g_name = '' OR p_g_name IS NULL) AND
     (old_g_desc = p_g_desc OR p_g_desc = '' OR p_g_desc IS NULL)
  THEN
    RETURN FALSE;
  END IF;
  UPDATE "group"
     SET "name" = COALESCE (NULLIF (trim (p_g_name), ''), old_g_name),
         "description" = COALESCE (NULLIF (trim (p_g_desc), ''), old_g_desc),
         "updated_at" = 'now ()'
   WHERE "group"."group_id" = p_group_id AND
         "group"."owner_id" = p_owner_id;
  GET DIAGNOSTICS n_affected_rows = ROW_COUNT;
  RETURN n_affected_rows;
  IF n_affected_rows > 0 THEN
    RETURN TRUE;
  ELSE
    RETURN FALSE;
  END IF;
END;
$$;

ALTER FUNCTION update_user
      OWNER TO "noda";
