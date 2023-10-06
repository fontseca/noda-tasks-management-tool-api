CREATE OR REPLACE PROCEDURE validate_rpp_and_page (
  INOUT p_rpp BIGINT, /* records per page  */
  INOUT p_page BIGINT
)
LANGUAGE 'plpgsql'
AS $$
DECLARE
  max_int_64 CONSTANT BIGINT := 9223372036854775807;
  max_valid_page_before_overflow BIGINT;
BEGIN
  /* Make sure we can retrieve `p_rpp' records in just
     `p_page' pages.  If not, then use the maximum value
     for a page.  */
  IF p_rpp <= 0 OR p_rpp IS NULL THEN
    p_rpp := 1;
  END IF;
  IF p_page <= 0 OR p_page IS NULL THEN
    p_page := 1;
  END IF;
  max_valid_page_before_overflow := (max_int_64 / p_rpp) - 1;
  IF p_page > max_valid_page_before_overflow THEN
    p_page := max_valid_page_before_overflow;
  END IF;
END;
$$;

ALTER PROCEDURE validate_rpp_and_page
       OWNER TO "noda";
