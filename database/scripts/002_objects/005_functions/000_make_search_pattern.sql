CREATE OR REPLACE FUNCTION make_search_pattern (
    IN p_keyword TEXT
)
RETURNS TEXT
LANGUAGE 'plpgsql'
AS $$
BEGIN
  IF p_keyword IS NULL OR p_keyword = '' THEN
    RETURN '';
  END IF;
  RETURN '(?=.*' || replace(trim(BOTH ' ' FROM lower(p_keyword)), ' ', '.*)(?=.*') || '.*).*';
END;
$$;

ALTER FUNCTION make_search_pattern
      OWNER TO "noda";
