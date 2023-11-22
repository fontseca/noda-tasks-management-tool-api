
CREATE OR REPLACE PROCEDURE assert_is_not_special_list (
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_list_id  "list"."list_id"%TYPE
)
  LANGUAGE 'plpgsql'
AS $$
DECLARE
  today_list_id "list"."list_id"%TYPE := get_today_list_id(p_owner_id);
  tomorrow_list_id "list"."list_id"%TYPE := get_tomorrow_list_id(p_owner_id);
BEGIN
  IF today_list_id = p_list_id OR tomorrow_list_id = p_list_id THEN
    RAISE EXCEPTION 'nonexistent list with ID "%"', p_list_id::TEXT
      USING HINT = 'Please check the given list ID.';
  END IF;
END;
$$;

ALTER PROCEDURE assert_is_not_special_list ("list"."owner_id"%TYPE,
  "list"."list_id"%TYPE)
  OWNER TO "noda";
