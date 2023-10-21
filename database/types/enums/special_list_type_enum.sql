DROP TYPE IF EXISTS "special_list_type_t";

CREATE TYPE "special_list_type_t"
    AS ENUM ('today',
             'tomorrow');

ALTER TYPE "special_list_type_t"
  OWNER TO "noda";
