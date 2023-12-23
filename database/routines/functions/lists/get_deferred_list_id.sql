CREATE OR REPLACE FUNCTION get_deferred_list_id
(
  IN p_user_id "list"."owner_id"%TYPE
)
RETURNS UUID
RETURNS NULL ON NULL INPUT
LANGUAGE 'sql'
AS $$
  SELECT "list_id"
    FROM "user_special_list"
   WHERE "user_id" = p_user_id
     AND "list_type" = 'deferred';
$$;

ALTER FUNCTION get_deferred_list_id ("list"."owner_id"%TYPE)
      OWNER TO "noda";
