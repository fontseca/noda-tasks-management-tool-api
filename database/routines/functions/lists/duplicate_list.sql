CREATE OR REPLACE FUNCTION duplicate_list
(
  IN p_owner_id "list"."owner_id"%TYPE,
  IN p_list_id  "list"."list_id"%TYPE
)
RETURNS "list"."list_id"%TYPE
LANGUAGE 'plpgsql'
AS $$
DECLARE
  current_list "list"%ROWTYPE;
  new_list_id "list"."list_id"%TYPE;
BEGIN
  CALL assert_user_exists (p_owner_id);
  CALL assert_is_not_special_list (p_owner_id, p_list_id);
  CALL assert_list_exists_somewhere (p_owner_id, p_list_id);
  SELECT *
    INTO current_list
    FROM "list"
   WHERE "owner_id" = p_owner_id AND
         "list_id" = p_list_id;
   new_list_id := make_list (p_owner_id,
                             current_list."group_id",
                             current_list."name",
                             current_list."description");
  PERFORM *
     FROM "task" t,
  LATERAL make_task (p_owner_id,
                     new_list_id,
                     ROW(t."title",
                         t."headline",
                         t."description",
                         t."priority",
                         t."status",
                         t."due_date",
                         t."remind_at"))
    WHERE "owner_id" = p_owner_id AND
          "list_id" = p_list_id;
  RETURN new_list_id;
END;
$$;

ALTER FUNCTION duplicate_list ("list"."owner_id"%TYPE,
                               "list"."list_id"%TYPE)
      OWNER TO "noda";
