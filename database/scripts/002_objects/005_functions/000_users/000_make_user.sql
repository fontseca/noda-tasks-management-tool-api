CREATE OR REPLACE FUNCTION make_user (
  IN p_first_name TEXT,
  IN p_middle_name TEXT,
  IN p_last_name TEXT,
  IN p_surname TEXT,
  IN p_email email_t,
  IN p_password TEXT
)
RETURNS UUID
LANGUAGE 'plpgsql'
AS $$
DECLARE
  new_user_id UUID;
BEGIN
  INSERT INTO "user"
              ("first_name", "middle_name", "last_name", "surname", "email", "password", "role_id")
       VALUES (p_first_name, p_middle_name, p_last_name, p_surname, p_email, p_password, '2')
    RETURNING "user_id"
         INTO new_user_id;
  RETURN new_user_id;
END;
$$;

ALTER FUNCTION make_user
      OWNER TO "noda";
