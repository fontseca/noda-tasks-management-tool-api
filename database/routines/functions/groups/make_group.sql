CREATE OR REPLACE FUNCTION make_group (
  IN p_owner_id "group"."owner_id"%TYPE,
  IN p_g_name   "group"."name"%TYPE,
  IN p_g_desc   "group"."description"%TYPE
)
RETURNS UUID
LANGUAGE 'plpgsql'
AS $$
DECLARE
  inserted_id UUID;
BEGIN
  CALL assert_user_exists (p_owner_id);
  INSERT INTO "group" ("owner_id", "name", "description")
       VALUES (p_owner_id, p_g_name, p_g_desc)
    RETURNING "group_id"
         INTO inserted_id;
  RETURN inserted_id;
END;
$$;

ALTER FUNCTION make_group
      OWNER TO "noda";
