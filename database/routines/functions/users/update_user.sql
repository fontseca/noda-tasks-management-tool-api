CREATE OR REPLACE FUNCTION update_user (
  IN p_user_id     "user"."user_id"%TYPE,
  IN p_first_name  "user"."first_name"%TYPE,
  IN p_middle_name "user"."middle_name"%TYPE,
  IN p_last_name   "user"."last_name"%TYPE,
  IN p_surname     "user"."surname"%TYPE,
  IN p_email       "user"."email"%TYPE,
  IN p_picture_url "user"."picture_url"%TYPE,
  IN p_password    "user"."password"%TYPE
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  rows_affected INT;
  old_first_name "user"."first_name"%TYPE;
  old_middle_name "user"."middle_name"%TYPE;
  old_last_name "user"."last_name"%TYPE;
  old_surname "user"."surname"%TYPE;
  old_email TEXT;
  old_picture_url "user"."picture_url"%TYPE;
  old_password "user"."password"%TYPE;
BEGIN
  CALL assert_user_exists (p_user_id);
  SELECT u."first_name",
         u."middle_name",
         u."last_name",
         u."surname",
         u."email",
         u."picture_url",
         u."password"
    INTO old_first_name,
         old_middle_name,
         old_last_name,
         old_surname,
         old_email,
         old_picture_url,
         old_password
    FROM "user" u
   WHERE u."user_id" = p_user_id;
  IF (old_first_name = p_first_name OR p_first_name = '' OR p_first_name IS NULL) AND
     (old_middle_name = p_middle_name OR p_middle_name = '' OR p_middle_name IS NULL) AND
     (old_last_name = p_last_name OR p_last_name = '' OR p_last_name IS NULL) AND
     (old_surname = p_surname OR p_surname = '' OR p_surname IS NULL) AND
     (old_email = p_email OR p_email = '' OR p_email IS NULL) AND
     (old_picture_url = p_picture_url OR p_picture_url = '' OR p_picture_url IS NULL) AND
     (old_password = p_password OR p_password = '' OR p_password IS NULL)
  THEN
    RETURN FALSE;
  END IF;
  UPDATE "user"
     SET "first_name" = COALESCE (NULLIF (trim (p_first_name), ''), old_first_name),
         "middle_name" = COALESCE (NULLIF (trim (p_middle_name), ''), old_middle_name),
         "last_name" = COALESCE (NULLIF (trim (p_last_name), ''), old_last_name),
         "surname" = COALESCE (NULLIF (trim (p_surname), ''), old_surname),
         "email" = COALESCE (NULLIF (trim (p_email), ''), old_email),
         "picture_url" = COALESCE (NULLIF (trim (p_picture_url), ''), old_picture_url),
         "password" = COALESCE (NULLIF (trim (p_password), ''), old_password),
         "updated_at" = 'now ()'
   WHERE "user"."user_id" = p_user_id;
  GET DIAGNOSTICS rows_affected = ROW_COUNT;
  RETURN rows_affected;
END;
$$;

ALTER FUNCTION update_user ("user"."user_id"%TYPE,
                            "user"."first_name"%TYPE,
                            "user"."middle_name"%TYPE,
                            "user"."last_name"%TYPE,
                            "user"."surname"%TYPE,
                            "user"."email"%TYPE,
                            "user"."picture_url"%TYPE,
                            "user"."password"%TYPE)
      OWNER TO "noda";
