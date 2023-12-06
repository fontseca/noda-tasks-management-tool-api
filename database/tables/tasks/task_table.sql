CREATE TABLE IF NOT EXISTS "task"
(
  "task_id"          UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4 (),
  "owner_id"         UUID NOT NULL REFERENCES "user" ("user_id") ON DELETE CASCADE,
  "list_id"          UUID NOT NULL REFERENCES "list" ("list_id") ON DELETE CASCADE,
  "position_in_list" pos_t NOT NULL,
  "title"            VARCHAR(100) NOT NULL,
  "headline"         VARCHAR DEFAULT NULL,
  "description"      TEXT DEFAULT NULL,
  "priority"         task_priority_t NOT NULL,
  "status"           task_status_t NOT NULL,
  "is_pinned"        BOOLEAN NOT NULL DEFAULT FALSE,
  "due_date"         TIMESTAMPTZ DEFAULT NULL,
  "remind_at"        TIMESTAMPTZ DEFAULT NULL,
  "completed_at"     TIMESTAMPTZ DEFAULT NULL,
  "created_at"       TIMESTAMPTZ NOT NULL DEFAULT now (),
  "updated_at"       TIMESTAMPTZ NOT NULL DEFAULT now ()
);

ALTER TABLE "task"
   OWNER TO "noda";

COMMENT ON TABLE "task"
              IS 'Manages individual tasks, including titles, descriptions, statuses, etc.';
