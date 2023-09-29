CREATE TABLE IF NOT EXISTS "role" (
  "role_id" SERIAL NOT NULL PRIMARY KEY,
  "name"    TEXT NOT NULL
);

ALTER TABLE "role"
   OWNER TO "noda";

COMMENT ON TABLE "role"
              IS 'A role for a user.';
