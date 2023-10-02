INSERT INTO "user_setting" ("user_id", "key", "value")
     SELECT "user_id",
            "key",
            "default_value"
       FROM "user",
            "predefined_user_setting"
   ORDER BY "user_id";
