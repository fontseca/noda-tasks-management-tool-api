CREATE OR REPLACE FUNCTION fetch_user_by_id (IN p_user_id "user"."user_id"%TYPE)
RETURNS "user"
LANGUAGE 'sql'
AS $$
  SELECT *
    FROM fetch_user_by ('user_id', p_user_id::TEXT)
$$;

ALTER FUNCTION fetch_user_by_id ("user"."user_id"%TYPE)
      OWNER TO "noda";
