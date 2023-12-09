CREATE TABLE IF NOT EXISTS "list"
(
  "list_id"     UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4 (),
  "owner_id"    UUID NOT NULL REFERENCES "user" ("user_id") ON DELETE CASCADE,
  "group_id"    UUID DEFAULT NULL REFERENCES "group" ("group_id") ON DELETE CASCADE,
  "name"        VARCHAR(64) NOT NULL,
  "description" VARCHAR(512) DEFAULT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT now (),
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT now ()
);

ALTER TABLE "list"
   OWNER TO "noda";

COMMENT ON TABLE "list"
              IS 'Organizes tasks under a single unit.';
