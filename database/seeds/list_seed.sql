INSERT INTO "list" ("owner_id", "name", "description")
     SELECT "user_id",
            'today@' || "user_id",
            "first_name" || ' ' || "last_name" || '''s today list'
       FROM "user"
  UNION ALL
     SELECT "user_id",
            'tomorrow@' || "user_id",
            "first_name" || ' ' || "last_name" || '''s tomorrow list'
       FROM "user"
   ORDER BY "user_id";
