CREATE OR REPLACE FUNCTION fetch_user_by (
  IN p_column TEXT,
  IN p_value  TEXT
)
RETURNS SETOF "user"
LANGUAGE 'plpgsql'
AS $$
BEGIN
  IF p_column IS NULL THEN
    p_column = '(NULL)';
  ELSE
    p_column := lower (trim (BOTH ' ' FROM p_column));
  END IF;
  IF p_column = 'user_id' THEN
    CALL assert_user_exists (p_value::UUID);
    RETURN QUERY
          SELECT *
            FROM "user"
           WHERE "user_id" = p_value::UUID;
  ELSIF p_column = 'email' THEN
  RETURN QUERY
    SELECT *
      FROM "user"
     WHERE "email" = p_value::email_t;
    IF NOT FOUND THEN
      IF p_value IS NULL THEN
        p_value := '(NULL)';
      END IF;
      RAISE EXCEPTION 'nonexistent user email: %', p_value
           USING HINT = 'Please check the given user email';
    END IF;
  ELSE
    RAISE EXCEPTION 'unexpected input for "p_column" parameter: "%"', p_column
       USING DETAIL = 'The supported values for "p_column" are "email" and "user_id"',
               HINT = 'Please check the given user email';
  END IF;
END;
$$;

ALTER FUNCTION fetch_user_by
      OWNER TO "noda";
