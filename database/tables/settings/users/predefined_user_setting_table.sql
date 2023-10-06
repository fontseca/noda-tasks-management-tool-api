CREATE TABLE IF NOT EXISTS "predefined_user_setting" (
  "key"           TEXT PRIMARY KEY,
  "default_value" JSON NULL,
  "description"   TEXT NULL
);

ALTER TABLE "predefined_user_setting"
   OWNER TO "noda";

COMMENT ON TABLE "predefined_user_setting"
              IS 'Default user-specific settings.';
