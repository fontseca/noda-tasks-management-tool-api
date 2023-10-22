CREATE OR REPLACE FUNCTION fetch_user_by_email (IN p_user_email TEXT)
RETURNS "user"
LANGUAGE 'sql'
AS $$
  SELECT *
    FROM fetch_user_by ('email', p_user_email)
$$;

ALTER FUNCTION fetch_user_by_email (TEXT)
      OWNER TO "noda";
