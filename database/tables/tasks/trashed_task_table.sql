CREATE TABLE IF NOT EXISTS "trashed_task"
                  AS TABLE "task";

               ALTER TABLE "trashed_task"
  ADD COLUMN IF NOT EXISTS "trashed_at" TIMESTAMPTZ NOT NULL DEFAULT now (),
  ADD COLUMN IF NOT EXISTS "destroy_at" TIMESTAMPTZ NOT NULL DEFAULT now () + INTERVAL '7d';

ALTER TABLE "trashed_task"
  OWNER TO "noda";
