CREATE OR REPLACE FUNCTION quote_meta (IN pattern TEXT)
RETURNS TEXT
LANGUAGE 'plpgsql'
AS $$
DECLARE
  special_characters CONSTANT TEXT[] := ARRAY[ '\\', '.', '*', '+', '?', '|', '(', ')', '[', ']', '{', '}', '^', '$' ];
  special_char TEXT;
BEGIN
    FOREACH special_char IN ARRAY special_characters
    LOOP
        pattern := replace (pattern, special_char, E'\\' || special_char);
    END LOOP;
    RETURN pattern;
END;
$$;

ALTER FUNCTION quote_meta (pattern TEXT)
      OWNER TO "noda";
