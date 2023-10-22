CREATE OR REPLACE FUNCTION check_setting_is_not_repeated()
RETURNS TRIGGER
LANGUAGE 'plpgsql'
AS $$
DECLARE
  n_settings INT;
BEGIN
  SELECT count (*)
    INTO n_settings
    FROM "user_setting" 
    WHERE "key" = NEW."key" AND
          "user_id" = NEW."user_id";
    IF n_settings >= 1 THEN
        RAISE EXCEPTION 'Key (key)=(%) is already set for user with ID ''%s''', NEW."key", NEW."user_id";
  END IF;
  RETURN NEW;
END;
$$;

ALTER FUNCTION check_setting_is_not_repeated()
      OWNER TO "noda";

    CREATE OR REPLACE TRIGGER check_setting_is_not_repeated
             BEFORE INSERT ON "user_setting"
FOR EACH ROW EXECUTE FUNCTION check_setting_is_not_repeated();
