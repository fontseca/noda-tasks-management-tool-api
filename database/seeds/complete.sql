
-- CREATE FUNCTION somefunc() RETURNS integer AS $$
-- << outerblock >>
-- DECLARE
--     quantity integer := 30;
-- BEGIN
--     RAISE NOTICE 'Quantity here is %', quantity;  -- Prints 30
--     quantity := 50;
--     --
--     -- Create a subblock
--     --
--     DECLARE
--         quantity integer := 80;
--     BEGIN
--         RAISE NOTICE 'Quantity here is %', quantity;  -- Prints 80
--         RAISE NOTICE 'Outer quantity here is %', outerblock.quantity;  -- Prints 50
--     END;

--     RAISE NOTICE 'Quantity here is %', quantity;  -- Prints 50

--     RETURN quantity;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE OR REPLACE FUNCTION baz()
-- RETURNS VOID
-- LANGUAGE 'plpgsql'
-- AS $$
-- BEGIN
--   RAISE NOTICE 'Heyyy';
-- END $$;



-- CREATE OR REPLACE FUNCTION foo() RETURNS INTEGER AS $$
--   SELECT 1 AS result;
-- $$ LANGUAGE ' sql';


-- CREATE FUNCTION one() RETURNS integer AS '
--     SELECT 1 AS result;
-- ' LANGUAGE SQL;

-- SELECT $$
--   Hello "World"
--   {}
--   BEGIN
--     SELECT * FROM "user"
--   END
-- $$;

-- CREATE OR REPLACE FUNCTION "make_user" (
--   IN "p_first_name" TEXT,
--   IN "p_middle_name" TEXT,
--   IN "p_last_name" TEXT,
--   IN "p_surname" TEXT,
--   IN "p_email" email_t,
--   IN "p_password" TEXT
-- )
-- RETURNS TABLE (
--   "id" UUID,
--   "role" INTEGER,
--   "first_name" TEXT,
--   "middle_name" TEXT,
--   "last_name" TEXT,
--   "surname" TEXT,
--   "picture_url" TEXT,
--   "email" TEXT,
--   "created_at" TIMESTAMPTZ,
--   "updated_at" TIMESTAMPTZ
-- )
-- AS $$
-- BEGIN
-- 	INSERT INTO "user"
-- 	            ("first_name", "middle_name", "last_name", "surname", "email", "password", "role_id")
--        VALUES ("p_first_name", "p_middle_name", "p_last_name", "p_surname", "p_email", "p_password", '2')
--     RETURNING "user_id" AS "id",
-- 		          "role_id" AS "role",
--               "first_name",
--               "middle_name",
--               "last_name",
--               "surname",
--               "picture_url",
--               "email",
--               "created_at",
--               "updated_at";
-- END;
-- $$ LANGUAGE 'plpgsql';


-- CREATE OR REPLACE FUNCTION add_one(IN val INTEGER DEFAULT 0)
--            RETURNS INTEGER
--           LANGUAGE PLPGSQL
--                     AS $$
-- BEGIN
--   RETURN val + 1;
-- END $$;

-- DROP FUNCTION IF EXISTS select_name;

-- CREATE OR REPLACE FUNCTION select_name()
--              RETURNS TABLE (
--               "first_name" TEXT,
--               "last_name" TEXT
--               )
--           LANGUAGE PLPGSQL
--                     AS $$
-- BEGIN
--   RETURN QUERY SELECT 'Jeremy', 'Fonseca';
--   RETURN QUERY SELECT 'Alexander', 'Blanco';
--   -- INSERT INTO "first_name" VALUES 'jeremy';
-- END $$;


-- SELECT * from select_name();

-- SELECT add_one(1) AS "Result", 'hey';

-- SELECT add_one(-1) AS "Result", 'hey';

-- -- select * from (
-- --   select user_setting_id, user_id, "key", "value", description, created_at, updated_at from user_setting inner join predfined_user_setting on user_setting.key = predfined_user_setting.key) where user_id = '01de206c-d562-402b-bff0-b38cef927807'

--     SELECT "us"."key",
--            "df"."description",
--            "us"."value",
--            "us"."created_at",
--            "us"."updated_at"
--       FROM "user_setting" "us"
-- INNER JOIN "predefined_user_setting" "df"
--         ON "us"."key" = "df"."key"
--      WHERE "us"."user_id" = 'a33c37df-e0c1-4746-bf09-57cda65bfd0c';

-- WITH "new_users" AS
-- (
--   INSERT INTO "user"
--               ("role_id", "first_name", "middle_name", "last_name",   "surname",   "picture_url",                                     "email",                    "password")
--        VALUES (1,         'Jeremy',     'Alexander',   'Fonseca',     'Blanco',    'http://dummyimage.com/235x100.png/5fa2dd/ffffff', 'f@mail.com',               '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Harmonie',   'Merla',       'Saberton',    'Shoulders', 'http://dummyimage.com/163x100.png/cc0000/ffffff', 'mshoulders1@shop-pro.jp',  '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Roscoe',     'Merry',       'Sibyllina',   'Dixson',    'http://dummyimage.com/145x100.png/cc0000/ffffff', 'mdixson2@typepad.com',     '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Nedda',      'Kristin',     'Lewin',       'Crispin',   'http://dummyimage.com/207x100.png/5fa2dd/ffffff', 'kcrispin3@alexa.com',      '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Noelyn',     'Muriel',      'De Few',      'Fewkes',    'http://dummyimage.com/165x100.png/5fa2dd/ffffff', 'mfewkes4@photobucket.com', '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Austin',     'Skyler',      'Kitchenside', 'Masson',    'http://dummyimage.com/153x100.png/dddddd/000000', 'smasson5@blogspot.com',    '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Shirlene',   'Illa',        'Staynes',     'MacAless',  'http://dummyimage.com/176x100.png/cc0000/ffffff', 'imacaless6@ameblo.jp',     '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Sherrie',    'Hamnet',      'Prestedge',   'Fackney',   'http://dummyimage.com/214x100.png/5fa2dd/ffffff', 'hfackney7@patch.com',      '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Jared',      'Catlaina',    'McFarlane',   'Craighill', 'http://dummyimage.com/212x100.png/ff4444/ffffff', 'ccraighill8@blogs.com',    '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC'),
--               (2,         'Eba',        'Raynard',     'Yakovl',      'Gurnett',   'http://dummyimage.com/182x100.png/dddddd/000000', 'rgurnett9@guardian.co.uk', '$2a$10$ODkl/7qJaSbpg025Ddhu5ewRDsKUOk8L2M6YYnzIW4N9t8mAyGRHC')
--     RETURNING "user_id", "first_name", "last_name"
-- )
-- INSERT INTO "list" ("owner_id", "name", "description")
--      SELECT "user_id",
--             'today',
--             "first_name" || ' ' || "last_name" || '''s today list'
--        FROM "new_users"
--   UNION ALL
--      SELECT "user_id",
--             'tomorrow',
--             "first_name" || ' ' || "last_name" || '''s tomorrow list'
--        FROM "new_users";


-- -- We can add extra values that will be selected as columns.
-- SELECT "user_id",
--        "first_name",
--        'message:',
--        'hello',
--        'jeremy' || ' ' || 'fonseca' AS "full_name"
--   FROM "user";

-- -- In this way, we can use this feature to do the following:
-- INSERT INTO "list" ("owner_id", "name", "description")
--      SELECT "user_id", 'Today', 'Tasks from today.'
--        FROM "user";

-- -- We can copy a table in this way:
-- CREATE TABLE IF NOT EXISTS "list_trash"
--                   AS TABLE "list"
--               WITH NO DATA;

-- ALTER TABLE "list_trash"
--  ADD COLUMN "trashed_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
--  ADD COLUMN "destroy_at" TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30d';

-- SELECT *
--   FROM "list_trash";

-- -- Move lists to trash.
-- WITH "removed_lists" AS
-- (
--   DELETE FROM "list"
--     RETURNING *
-- )
-- INSERT INTO "list_trash"
--      SELECT *
--        FROM "removed_lists";

-- -- Create a bunch of users and create their today lists.
-- WITH "new_users" AS
-- (
--   INSERT INTO "user" ("first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password")
--        VALUES ('Briano', 'Geoff', 'Cakebread', 'Frie', 'http://dummyimage.com/235x100.png/5fa2dd/ffffff', 'gfrie0@addtoany.com', '$2a$04$5y2S0QM/Iz49Tc8sCOFtQOop3ipv0MkhqKpdOItoZRiUOsjMjbqvq'),
--               ('Harmonie', 'Merla', 'Saberton', 'Shoulders', 'http://dummyimage.com/163x100.png/cc0000/ffffff', 'mshoulders1@shop-pro.jp', '$2a$04$Rbdr8lc4I5V0uT5gffyUKOuH.32UDuPH436EaCA9MS5sbikN57wr.'),
--               ('Roscoe', 'Merry', 'Sibyllina', 'Dixson', 'http://dummyimage.com/145x100.png/cc0000/ffffff', 'mdixson2@typepad.com', '$2a$04$Xk1sjdPgeyuTVTP2HqXAlelWBta/84TSfuSl9J04XYon.uDcIUhza'),
--               ('Nedda', 'Kristin', 'Lewin', 'Crispin', 'http://dummyimage.com/207x100.png/5fa2dd/ffffff', 'kcrispin3@alexa.com', '$2a$04$FdVWjWbfN0bNVpmY4njLy.OMpnWVJm99waU8upRjINXOdWfqheZia'),
--               ('Noelyn', 'Muriel', 'De Few', 'Fewkes', 'http://dummyimage.com/165x100.png/5fa2dd/ffffff', 'mfewkes4@photobucket.com', '$2a$04$M8mgK.kDBpAfkuVwN/thg.QP7KMB9ZU0P3SNNFPmQORNyGHK8dQR2'),
--               ('Austin', 'Skyler', 'Kitchenside', 'Masson', 'http://dummyimage.com/153x100.png/dddddd/000000', 'smasson5@blogspot.com', '$2a$04$nFHMVqFlVJNkW1qP7LuFMef2jcJ0Pka5PEX70nOq0CdIZqotv/8mK'),
--               ('Shirlene', 'Illa', 'Staynes', 'MacAless', 'http://dummyimage.com/176x100.png/cc0000/ffffff', 'imacaless6@ameblo.jp', '$2a$04$vFBFYZjg7dPI6pgywm.uA.6HlB/cjUqNIxTpnn4Oy0jFQRfH26PwO'),
--               ('Sherrie', 'Hamnet', 'Prestedge', 'Fackney', 'http://dummyimage.com/214x100.png/5fa2dd/ffffff', 'hfackney7@patch.com', '$2a$04$PmfX.ivAlUDkgP6I4RT7X.S5vKs2.d6z/PynV34AHbqIk/RjwopZO'),
--               ('Jared', 'Catlaina', 'McFarlane', 'Craighill', 'http://dummyimage.com/212x100.png/ff4444/ffffff', 'ccraighill8@blogs.com', '$2a$04$yW9sAf/2ddeDiaRQTxnmWOC362mUhXwN0SV96euO4DqMY9JLctLR2'),
--               ('Eba', 'Raynard', 'Yakovl', 'Gurnett', 'http://dummyimage.com/182x100.png/dddddd/000000', 'rgurnett9@guardian.co.uk', '$2a$04$7yRXvsngK7c.EyizjINieutuQh1.fIoj8FaNP4mdDQ8/4yKMUduz.')
--     RETURNING "user_id", "last_name"
-- )
-- INSERT INTO "list" ("owner_id", "name", "description")
--      SELECT "user_id",
--             'today',
--             "last_name" || '''s today tasks'
--        FROM "new_users";

-- -- Create a single user and create its today list.
-- WITH "new_user" AS
-- (
--   INSERT INTO "user" ("first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password") VALUES ('Briano', 'Geoff', 'Cakebread', 'Frie', 'http://dummyimage.com/235x100.png/5fa2dd/ffffff', 'gfrie0@addtoany.com', '$2a$04$5y2S0QM/Iz49Tc8sCOFtQOop3ipv0MkhqKpdOItoZRiUOsjMjbqvq')
--     RETURNING *
-- )
-- INSERT INTO "list" ("owner_id", "name", "description")
--      SELECT "user_id", 'Today', 'today tasks for '||"last_name"
--        FROM "new_user";

-- -----------

-- -- insert into recipient(name, address) values(name, address);
-- -- insert into orders(recipient_id, total_price, total_quantity) values(recipient_id, 2000, 20);
-- -- insert into items(order_id, item, price, quantity, total) values(order_id, item1, 230, 2, 260);
-- -- insert into items(order_id, item, price, quantity, total) values(order_id, item2, 500, 2, 1000);

-- CREATE TABLE "recipient"
-- (
--   recipient_id SERIAL PRIMARY KEY,
--   name text,
--   address text
-- );

-- CREATE TABLE "order"
-- (
--   order_id SERIAL PRIMARY KEY,
--   recipient_id SERIAL,
--   total_price int,
--   total_quantity int
-- );

-- CREATE TABLE "items"
-- (
--   item_id SERIAL PRIMARY KEY,
--   order_id SERIAL,
--   item text,
--   price int,
--   quantity int,
--   total int
-- )

-- WITH
--   "new_recipient" AS
--   (
--     INSERT INTO "recipient" ("name", "address")
--         VALUES ('Alexander', 'Chinandega')
--       RETURNING "recipient_id"
--   ),
--   "new_order" AS
--   (
--     INSERT INTO "order" ("recipient_id", "total_price", "total_quantity")
--          VALUES ((SELECT * FROM "new_recipient"), 2000, 20)
--       RETURNING "order_id"
--   )
--   INSERT INTO "items" ("order_id", "item", "price", "quantity", "total")
--        VALUES ((SELECT "order_id" FROM "new_order"), 'item1', 230, 2, 260), 
--               ((SELECT "order_id" FROM "new_order"), 'item2', 500, 2, 1000);

-- SELECT *
--   FROM
--   (
--     SELECT *
--       FROM "list"
--   ) AS "users"
-- UNION
-- SELECT *
--   FROM "list_trash";



-- DO $$
-- DECLARE 
--     user_row RECORD;
-- BEGIN
--   FOR user_row IN SELECT "user_id" FROM "user" LOOP
--     -- Access the "first_name" column of the current row using user_row.first_name
--     -- SELECT * FROM "task" WHERE "owner_id" = user_row.user_id;
--     RAISE NOTICE 'First Name: %', user_row.user_id;
--   END LOOP;
-- END $$;
