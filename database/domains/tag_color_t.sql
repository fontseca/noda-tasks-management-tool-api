DROP DOMAIN IF EXISTS "tag_color_t";

CREATE DOMAIN "tag_color_t"
           AS VARCHAR NOT NULL
        CHECK (VALUE ~ '^([A-F0-9]{6})$' AND Length(VALUE) = 6);

ALTER DOMAIN "tag_color_t"
    OWNER TO "noda";

COMMENT ON DOMAIN "tag_color_t"
               IS 'Hexadecimal number used to define valid colors. It does not requires the #.';

