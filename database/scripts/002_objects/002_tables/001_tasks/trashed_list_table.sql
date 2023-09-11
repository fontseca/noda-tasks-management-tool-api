CREATE TABLE IF NOT EXISTS "trashed_list"
                  AS TABLE "list";

ALTER TABLE "trashed_list"
 ADD COLUMN "trashed_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
 ADD COLUMN "will_destroy_at" TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30d';

ALTER TABLE "trashed_list"
   OWNER TO "noda";
