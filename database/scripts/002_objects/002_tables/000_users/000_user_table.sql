CREATE TABLE IF NOT EXISTS "user"
(
  "user_id"     UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
  "first_name"  VARCHAR(50) NOT NULL,
  "middle_name" VARCHAR(50) DEFAULT NULL,
  "last_name"   VARCHAR(50) DEFAULT NULL,
  "surname"     VARCHAR(50) DEFAULT NULL,
  "picture_url" VARCHAR DEFAULT NULL,
  "email"       email_t UNIQUE,
  "password"    VARCHAR NOT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT now(),
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE "user"
   OWNER TO "noda";

COMMENT ON TABLE "user"
              IS 'Represents system users with their personal information and account details.';
