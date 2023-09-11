DROP TYPE IF EXISTS "task_priority_t";

CREATE TYPE "task_priority_t"
    AS ENUM
    (
      'high',
      'medium',
      'low'
    );

ALTER TYPE "task_priority_t"
  OWNER TO "noda";

COMMENT ON TYPE "task_priority_t"
             IS 'Defines the priority levels that can be assigned to tasks.';
