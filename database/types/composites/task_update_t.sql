DROP TYPE IF EXISTS "task_update_t";

CREATE TYPE "task_update_t" AS
(
  "title"       VARCHAR(128),
  "headline"    VARCHAR(64),
  "description" VARCHAR(512)
);

ALTER TYPE "task_update_t"
  OWNER TO "noda";

COMMENT ON TYPE "task_update_t"
             IS 'Represents the specifications for updating a task.';
