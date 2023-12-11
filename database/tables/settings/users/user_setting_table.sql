CREATE TABLE IF NOT EXISTS "user_setting"
(
  "user_setting_id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  "user_id"         UUID NOT NULL REFERENCES "user" ("user_id"),
  "key"             VARCHAR(50) NOT NULL REFERENCES "predefined_user_setting" ("key") ON DELETE CASCADE,
  "value"           JSON NOT NULL,
  "created_at"      TIMESTAMPTZ NOT NULL DEFAULT now (),
  "updated_at"      TIMESTAMPTZ NOT NULL DEFAULT now ()
);

ALTER TABLE "user_setting"
   OWNER TO "noda";

COMMENT ON TABLE "user_setting"
              IS 'User-specific settings represented as key-value pairs.';
