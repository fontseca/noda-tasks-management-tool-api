CREATE TABLE IF NOT EXISTS "blocked_user"
                  AS TABLE "user";

ALTER TABLE "blocked_user"
 ADD COLUMN "reason"     VARCHAR(512) DEFAULT NULL,
 ADD COLUMN "blocked_at" TIMESTAMPTZ NOT NULL DEFAULT now (),
 ADD COLUMN "blocked_by" UUID REFERENCES "user" ("user_id");

ALTER TABLE "blocked_user"
   OWNER TO "noda";
