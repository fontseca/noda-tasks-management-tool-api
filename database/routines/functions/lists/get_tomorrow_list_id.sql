CREATE OR REPLACE FUNCTION get_tomorrow_list_id (IN p_user_id "list"."owner_id"%TYPE)
RETURNS UUID
LANGUAGE 'sql'
AS $$
  SELECT "list_id"
    FROM "user_special_list"
   WHERE "user_id" = p_user_id AND
         "list_type" = 'tomorrow';
$$;

ALTER FUNCTION get_tomorrow_list_id ("list"."owner_id"%TYPE)
      OWNER TO "noda";
