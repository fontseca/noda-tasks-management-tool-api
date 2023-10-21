CREATE OR REPLACE FUNCTION make_today_list (IN p_owner_id "list"."owner_id"%TYPE)
RETURNS "list"."list_id"%TYPE
LANGUAGE 'plpgsql'
AS $$
DECLARE
  today_list_id "list"."list_id"%TYPE;
  existent_list_id "list"."list_id"%TYPE;
BEGIN
  CALL assert_user_exists (p_owner_id);
  SELECT get_today_list_id (p_owner_id)
    INTO existent_list_id;
  IF existent_list_id IS NOT NULL THEN
    RAISE EXCEPTION 'today list already exists for user with ID "%"', p_owner_id
              USING HINT = 'Function "make_today_list" should be invoked once per user.';
  END IF;
  INSERT INTO "list" ("owner_id", "group_id", "name", "description")
       SELECT p_owner_id,
              NULL,
              '___today___',
              concat (u."first_name", ' ', u."last_name", '''s today list')
         FROM "user" u
        WHERE  u."user_id" = p_owner_id
    RETURNING "list_id"
         INTO today_list_id;
  INSERT INTO "user_special_list" ("user_id", "list_id", "list_type")
       VALUES (p_owner_id, today_list_id, 'today');
  RETURN today_list_id;
END;
$$;

ALTER FUNCTION make_today_list ("list"."owner_id"%TYPE)
      OWNER TO "noda";
