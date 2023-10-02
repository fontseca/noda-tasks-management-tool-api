INSERT INTO "predefined_user_setting"
            ("key",                  "default_value", "description")
     VALUES ('language',             '"en_US.UTF-8"', 'The default user''s profile language.'),
            ('exported_data_format',         '"pdf"', 'The format to export the users'' tasks.'),
            ('show_task_progress',            'true', 'Shows the progress of a task based on accomplished steps.'),
            ('tz',                            'null', 'The users'' time zone.');
