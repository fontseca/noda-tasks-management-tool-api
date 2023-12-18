DROP TYPE IF EXISTS "special_list_type_t";

CREATE TYPE "special_list_type_t"
    AS ENUM ('today',
             'tomorrow',
             'deferred');

ALTER TYPE "special_list_type_t"
  OWNER TO "noda";
