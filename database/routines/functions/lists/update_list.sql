CREATE OR REPLACE FUNCTION update_list (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_group_id "list"."group_id"%TYPE,
  IN p_list_id  "list"."list_id"%TYPE,
  IN p_l_name   "list"."name"%TYPE,
  IN p_l_desc   "list"."description"%TYPE
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_affected_rows INT;
  old_l_name "list"."name"%TYPE;
  old_l_desc "list"."description"%TYPE;
  is_scattered_list CONSTANT BOOLEAN := p_group_id IS NULL;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_is_not_special_list (p_owner_id, p_list_id);
  CALL assert_list_exists (p_owner_id, p_group_id, p_list_id);
  SELECT l."name",
         l."description"
    INTO old_l_name,
         old_l_desc
    FROM "list" l
   WHERE l."owner_id" = p_owner_id AND
         CASE WHEN is_scattered_list
              THEN TRUE
              ELSE "group_id" = p_group_id
          END AND
         l."list_id" = p_list_id;
  p_l_name := NULLIF (trim (BOTH ' ' FROM p_l_name), '');
  p_l_desc := trim (BOTH ' ' FROM p_l_desc);
  IF (p_l_name IS NULL OR old_l_name = p_l_name) AND
     (p_l_desc IS NULL OR old_l_desc = p_l_desc)
  THEN
    RETURN FALSE;
  END IF;
  UPDATE "list"
     SET "name" = COALESCE (p_l_name, old_l_name),
         "description" = COALESCE (p_l_desc, old_l_desc),
         "updated_at" = 'now ()'
   WHERE "owner_id" = p_owner_id AND
         CASE WHEN is_scattered_list
              THEN "group_id" IS NULL
              ELSE "group_id" = p_group_id
          END AND
         "list_id" = p_list_id;
  GET DIAGNOSTICS n_affected_rows = ROW_COUNT;
  RETURN n_affected_rows;
END;
$$;

ALTER FUNCTION update_list ("list"."owner_id"%TYPE,
                            "list"."group_id"%TYPE,
                            "list"."list_id"%TYPE,
                            "list"."name"%TYPE,
                            "list"."description"%TYPE)
      OWNER TO "noda";
