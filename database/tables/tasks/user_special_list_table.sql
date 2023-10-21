DROP TABLE "user_special_list";

CREATE TABLE IF NOT EXISTS "user_special_list" (
  "user_special_list_id" UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4 (),
  "user_id"              UUID REFERENCES "user" ("user_id"),
  "list_id"              UUID REFERENCES "list" ("list_id"),
  "list_type"            special_list_type_t NOT NULL
);

ALTER TABLE "user_special_list"
   OWNER TO "noda";
