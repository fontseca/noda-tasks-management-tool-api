CREATE TABLE IF NOT EXISTS "attachment"
(
  "attachment_id" UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4 (),
  "owner_id"      UUID NOT NULL REFERENCES "user" ("user_id") ON DELETE CASCADE,
  "task_id"       UUID NOT NULL REFERENCES "task" ("task_id") ON DELETE CASCADE,
  "file_name"     TEXT NOT NULL,
  "file_url"      TEXT NOT NULL,
  "created_at"    TIMESTAMPTZ NOT NULL DEFAULT now (),
  "updated_at"    TIMESTAMPTZ NOT NULL DEFAULT now ()
);

ALTER TABLE "attachment"
   OWNER TO "noda";

COMMENT ON TABLE "attachment"
              IS 'Stores files associated with tasks.';
