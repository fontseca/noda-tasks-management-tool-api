CREATE OR REPLACE FUNCTION delete_user_hardly (IN p_user_id "user"."user_id"%TYPE)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_affected_rows INT;
BEGIN
  CALL assert_user_exists (p_user_id);
  DELETE FROM "user"
        WHERE "user_id" = p_user_id;
  GET DIAGNOSTICS n_affected_rows := ROW_COUNT;
  RETURN n_affected_rows;
END;
$$;

ALTER FUNCTION delete_user_hardly ("user"."user_id"%TYPE)
      OWNER TO "noda";
