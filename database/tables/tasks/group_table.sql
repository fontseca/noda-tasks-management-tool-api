CREATE TABLE IF NOT EXISTS "group"
(
  "group_id"    UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4 (),
  "owner_id"    UUID NOT NULL REFERENCES "user" ("user_id") ON DELETE CASCADE,
  "name"        VARCHAR(50) NOT NULL UNIQUE,
  "description" TEXT DEFAULT NULL,
  "is_archived" BOOLEAN NOT NULL DEFAULT FALSE,
  "archived_at" TIMESTAMPTZ DEFAULT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT now (),
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT now ()
);

ALTER TABLE "group"
   OWNER TO "noda";

COMMENT ON TABLE "group"
              IS 'Gather together lists.';
