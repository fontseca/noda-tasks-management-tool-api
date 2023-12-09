CREATE TABLE IF NOT EXISTS "predefined_user_setting" (
  "key"           VARCHAR(50) PRIMARY KEY,
  "default_value" JSON NULL,
  "description"   VARCHAR(512) NULL
);

ALTER TABLE "predefined_user_setting"
   OWNER TO "noda";

COMMENT ON TABLE "predefined_user_setting"
              IS 'Default user-specific settings.';
