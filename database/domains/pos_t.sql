DROP DOMAIN IF EXISTS "pos_t";

CREATE DOMAIN "pos_t"
           AS INTEGER NOT NULL DEFAULT 0
        CHECK (VALUE >= 0);

ALTER DOMAIN "pos_t"
    OWNER TO "noda";

COMMENT ON DOMAIN "pos_t"
               IS 'Represents a positional value.';
