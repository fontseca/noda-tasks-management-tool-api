CREATE TABLE IF NOT EXISTS "step"
(
  "step_id"      UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4 (),
  "task_id"      UUID NOT NULL REFERENCES "task" ("task_id"),
  "order"        pos_t NOT NULL UNIQUE,
  "description"  VARCHAR(512) DEFAULT NULL,
  "completed_at" TIMESTAMPTZ DEFAULT NULL,
  "created_at"   TIMESTAMPTZ NOT NULL DEFAULT now (),
  "updated_at"   TIMESTAMPTZ NOT NULL DEFAULT now ()
);

ALTER TABLE "step"
   OWNER TO "noda";

COMMENT ON TABLE "step"
              IS 'Logical steps to follow to complete a task.';
