DROP DOMAIN IF EXISTS "email_t";

CREATE DOMAIN "email_t"
           AS VARCHAR(240) NOT NULL
        CHECK (VALUE ~ '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$');

ALTER DOMAIN "email_t"
    OWNER TO "noda";

COMMENT ON DOMAIN "email_t"
               IS 'Ensures a string is a valid email format.';
