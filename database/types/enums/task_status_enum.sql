DROP TYPE IF EXISTS "task_status_t";

CREATE TYPE "task_status_t"
    AS ENUM ('finished',
             'in progress',
             'unfinished',
             'decayed');

ALTER TYPE "task_status_t"
  OWNER TO "noda";

COMMENT ON TYPE "task_status_t"
             IS 'Represents the different status levels that a task can have within the system.';
