CREATE OR REPLACE FUNCTION make_list (
  IN p_owner_id "list"."list_id"%TYPE,
  IN p_group_id "list"."group_id"%TYPE,
  IN p_l_name   "list"."name"%TYPE,
  IN p_l_desc   "list"."description"%TYPE
)
RETURNS "list"."list_id"%TYPE
LANGUAGE 'plpgsql'
AS $$
DECLARE
  is_scattered_list BOOLEAN := p_group_id IS NULL;
  inserted_list_id "list"."list_id"%TYPE;
  similar_names_count INT;
BEGIN
  CALL assert_user_exists (p_owner_id);
  IF is_scattered_list IS FALSE THEN
    CALL assert_group_exists (p_owner_id, p_group_id);
  END IF;
  p_l_name := NULLIF (trim (BOTH ' ' FROM p_l_name), '');
  p_l_desc := NULLIF (trim(BOTH ' ' FROM p_l_desc), '');
  SELECT count (*)
    INTO similar_names_count
    FROM "list" l
   WHERE CASE WHEN is_scattered_list THEN TRUE ELSE l."group_id" = p_group_id END AND
         regexp_count (l."name", '^' || quote_meta (p_l_name) || '(?: \(\d+\))?$') = 1;
  IF similar_names_count > 0 THEN
    p_l_name := concat (p_l_name, ' (', similar_names_count, ')');
  END IF;
  INSERT INTO "list" ("owner_id", "group_id", "name", "description")
       VALUES (p_owner_id,
               CASE WHEN is_scattered_list THEN NULL ELSE p_group_id END,
               p_l_name,
               p_l_desc)
    RETURNING "list_id"
         INTO inserted_list_id;
  RETURN inserted_list_id;
END;
$$;

ALTER FUNCTION make_list ("list"."list_id"%TYPE,
                          "list"."group_id"%TYPE,
                          "list"."name"%TYPE,
                          "list"."description"%TYPE)
      OWNER TO "noda";
