CREATE OR REPLACE FUNCTION make_deferred_list
(
  IN p_owner_id "list"."list_id"%TYPE
)
RETURNS "list"."list_id"%TYPE
RETURNS NULL ON NULL INPUT
LANGUAGE 'plpgsql'
AS $$
DECLARE
  deferred_list_id "list"."list_id"%TYPE;
  existent_list_id "list"."list_id"%TYPE;
BEGIN
  CALL assert_user_exists (p_owner_id);
  existent_list_id := get_deferred_list_id (p_owner_id);
  IF existent_list_id IS NOT NULL THEN
    RAISE EXCEPTION 'deferred list already exists for user with ID "%"', p_owner_id
         USING HINT = 'Function "make_deferred_list" should be invoked once per user.';
  END IF;
  INSERT INTO "list" ("owner_id", "group_id", "name", "description")
       SELECT p_owner_id,
              NULL,
              '___deferred___',
              concat (u."first_name", ' ',
                      u."last_name",
                      '''s deferred list')
         FROM "user" u
        WHERE  u."user_id" = p_owner_id
    RETURNING "list_id"
         INTO deferred_list_id;
  INSERT INTO "user_special_list" ("user_id", "list_id", "list_type")
       VALUES (p_owner_id, deferred_list_id, 'deferred');
  RETURN deferred_list_id;
END;
$$;

ALTER FUNCTION make_deferred_list ("list"."list_id"%TYPE)
      OWNER TO "noda";
