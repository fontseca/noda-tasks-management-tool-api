CREATE OR REPLACE FUNCTION unblock_user (IN p_user_id "user"."user_id"%TYPE)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  affected_rows INTEGER;
BEGIN
  WITH "user_to_unblock" AS
  (
    DELETE FROM "blocked_user"
          WHERE "user_id" = p_user_id
      RETURNING *
  )
  INSERT INTO "user" ("user_id",
                      "role_id",
                      "first_name",
                      "middle_name",
                      "last_name",
                      "surname",
                      "picture_url",
                      "email",
                      "password",
                      "created_at",
                      "updated_at")
       SELECT "user_id",
              "role_id",
              "first_name",
              "middle_name",
              "last_name",
              "surname",
              "picture_url",
              "email",
              "password",
              "created_at",
              "updated_at"
         FROM "user_to_unblock";
  GET DIAGNOSTICS affected_rows = ROW_COUNT;
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION unblock_user ("user"."user_id"%TYPE)
      OWNER TO "noda";
