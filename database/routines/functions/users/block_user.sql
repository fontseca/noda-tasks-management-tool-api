CREATE OR REPLACE FUNCTION block_user
(
  IN p_user_id    "user"."user_id"%TYPE,
  IN p_blocked_by "user"."user_id"%TYPE,
  IN p_reason     VARCHAR(512)
)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
DECLARE
  affected_rows INTEGER;
BEGIN
  CALL assert_user_exists (p_user_id);
  WITH "user_to_block" AS
  (
    DELETE FROM "user"
          WHERE "user_id" = p_user_id
      RETURNING *
  )
  INSERT INTO "blocked_user" ("user_id",
                              "role_id",
                              "first_name",
                              "middle_name",
                              "last_name",
                              "surname",
                              "picture_url",
                              "email",
                              "password",
                              "created_at",
                              "updated_at",
                              "reason",
                              "blocked_by")
       SELECT b."user_id",
              b."role_id",
              b."first_name",
              b."middle_name",
              b."last_name",
              b."surname",
              b."picture_url",
              b."email",
              b."password",
              b."created_at",
              b."updated_at",
              COALESCE (p_reason, 'unknown'),
              p_blocked_by
         FROM "user_to_block" b;
  GET DIAGNOSTICS affected_rows = ROW_COUNT;
  RETURN affected_rows;
END;
$$;

ALTER FUNCTION block_user ("user"."user_id"%TYPE,
                           "user"."user_id"%TYPE,
                           VARCHAR(512))
      OWNER TO "noda";

CREATE OR REPLACE FUNCTION block_user (IN p_user_id "user"."user_id"%TYPE)
RETURNS BOOLEAN
LANGUAGE 'plpgsql'
AS $$
BEGIN
  RETURN block_user(p_user_id, NULL, NULL);
END;
$$;

ALTER FUNCTION block_user ("user"."user_id"%TYPE)
      OWNER TO "noda";
