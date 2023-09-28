-- We can add extra values that will be selected as columns.
SELECT "user_id",
       "first_name",
       'message:',
       'hello',
       'jeremy' || ' ' || 'fonseca' AS "full_name"
  FROM "user";

-- In this way, we can use this feature to do the following:
INSERT INTO "list" ("owner_id", "name", "description")
     SELECT "user_id", 'Today', 'Tasks from today.'
       FROM "user";

-- We can copy a table in this way:
CREATE TABLE IF NOT EXISTS "list_trash"
                  AS TABLE "list"
              WITH NO DATA;

ALTER TABLE "list_trash"
 ADD COLUMN "trashed_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
 ADD COLUMN "destroy_at" TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30d';

SELECT *
  FROM "list_trash";

-- Move lists to trash.
WITH "removed_lists" AS
(
  DELETE FROM "list"
    RETURNING *
)
INSERT INTO "list_trash"
     SELECT *
       FROM "removed_lists";

-- Create a bunch of users and create their today lists.
WITH "new_users" AS
(
  INSERT INTO "user" ("first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password")
       VALUES ('Briano', 'Geoff', 'Cakebread', 'Frie', 'http://dummyimage.com/235x100.png/5fa2dd/ffffff', 'gfrie0@addtoany.com', '$2a$04$5y2S0QM/Iz49Tc8sCOFtQOop3ipv0MkhqKpdOItoZRiUOsjMjbqvq'),
              ('Harmonie', 'Merla', 'Saberton', 'Shoulders', 'http://dummyimage.com/163x100.png/cc0000/ffffff', 'mshoulders1@shop-pro.jp', '$2a$04$Rbdr8lc4I5V0uT5gffyUKOuH.32UDuPH436EaCA9MS5sbikN57wr.'),
              ('Roscoe', 'Merry', 'Sibyllina', 'Dixson', 'http://dummyimage.com/145x100.png/cc0000/ffffff', 'mdixson2@typepad.com', '$2a$04$Xk1sjdPgeyuTVTP2HqXAlelWBta/84TSfuSl9J04XYon.uDcIUhza'),
              ('Nedda', 'Kristin', 'Lewin', 'Crispin', 'http://dummyimage.com/207x100.png/5fa2dd/ffffff', 'kcrispin3@alexa.com', '$2a$04$FdVWjWbfN0bNVpmY4njLy.OMpnWVJm99waU8upRjINXOdWfqheZia'),
              ('Noelyn', 'Muriel', 'De Few', 'Fewkes', 'http://dummyimage.com/165x100.png/5fa2dd/ffffff', 'mfewkes4@photobucket.com', '$2a$04$M8mgK.kDBpAfkuVwN/thg.QP7KMB9ZU0P3SNNFPmQORNyGHK8dQR2'),
              ('Austin', 'Skyler', 'Kitchenside', 'Masson', 'http://dummyimage.com/153x100.png/dddddd/000000', 'smasson5@blogspot.com', '$2a$04$nFHMVqFlVJNkW1qP7LuFMef2jcJ0Pka5PEX70nOq0CdIZqotv/8mK'),
              ('Shirlene', 'Illa', 'Staynes', 'MacAless', 'http://dummyimage.com/176x100.png/cc0000/ffffff', 'imacaless6@ameblo.jp', '$2a$04$vFBFYZjg7dPI6pgywm.uA.6HlB/cjUqNIxTpnn4Oy0jFQRfH26PwO'),
              ('Sherrie', 'Hamnet', 'Prestedge', 'Fackney', 'http://dummyimage.com/214x100.png/5fa2dd/ffffff', 'hfackney7@patch.com', '$2a$04$PmfX.ivAlUDkgP6I4RT7X.S5vKs2.d6z/PynV34AHbqIk/RjwopZO'),
              ('Jared', 'Catlaina', 'McFarlane', 'Craighill', 'http://dummyimage.com/212x100.png/ff4444/ffffff', 'ccraighill8@blogs.com', '$2a$04$yW9sAf/2ddeDiaRQTxnmWOC362mUhXwN0SV96euO4DqMY9JLctLR2'),
              ('Eba', 'Raynard', 'Yakovl', 'Gurnett', 'http://dummyimage.com/182x100.png/dddddd/000000', 'rgurnett9@guardian.co.uk', '$2a$04$7yRXvsngK7c.EyizjINieutuQh1.fIoj8FaNP4mdDQ8/4yKMUduz.')
    RETURNING "user_id", "last_name"
)
INSERT INTO "list" ("owner_id", "name", "description")
     SELECT "user_id",
            'today',
            "last_name" || '''s today tasks'
       FROM "new_users";

-- Create a single user and create its today list.
WITH "new_user" AS
(
  INSERT INTO "user" ("first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password") VALUES ('Briano', 'Geoff', 'Cakebread', 'Frie', 'http://dummyimage.com/235x100.png/5fa2dd/ffffff', 'gfrie0@addtoany.com', '$2a$04$5y2S0QM/Iz49Tc8sCOFtQOop3ipv0MkhqKpdOItoZRiUOsjMjbqvq')
    RETURNING *
)
INSERT INTO "list" ("owner_id", "name", "description")
     SELECT "user_id", 'Today', 'today tasks for '||"last_name"
       FROM "new_user";

-----------

-- insert into recipient(name, address) values(name, address);
-- insert into orders(recipient_id, total_price, total_quantity) values(recipient_id, 2000, 20);
-- insert into items(order_id, item, price, quantity, total) values(order_id, item1, 230, 2, 260);
-- insert into items(order_id, item, price, quantity, total) values(order_id, item2, 500, 2, 1000);

CREATE TABLE "recipient"
(
  recipient_id SERIAL PRIMARY KEY,
  name text,
  address text
);

CREATE TABLE "order"
(
  order_id SERIAL PRIMARY KEY,
  recipient_id SERIAL,
  total_price int,
  total_quantity int
);

CREATE TABLE "items"
(
  item_id SERIAL PRIMARY KEY,
  order_id SERIAL,
  item text,
  price int,
  quantity int,
  total int
)

WITH
  "new_recipient" AS
  (
    INSERT INTO "recipient" ("name", "address")
        VALUES ('Alexander', 'Chinandega')
      RETURNING "recipient_id"
  ),
  "new_order" AS
  (
    INSERT INTO "order" ("recipient_id", "total_price", "total_quantity")
         VALUES ((SELECT * FROM "new_recipient"), 2000, 20)
      RETURNING "order_id"
  )
  INSERT INTO "items" ("order_id", "item", "price", "quantity", "total")
       VALUES ((SELECT "order_id" FROM "new_order"), 'item1', 230, 2, 260), 
              ((SELECT "order_id" FROM "new_order"), 'item2', 500, 2, 1000);

SELECT *
  FROM
  (
    SELECT *
      FROM "list"
  ) AS "users"
UNION
SELECT *
  FROM "list_trash";



DO $$
DECLARE 
    user_row RECORD;
BEGIN
  FOR user_row IN SELECT "user_id" FROM "user" LOOP
    -- Access the "first_name" column of the current row using user_row.first_name
    -- SELECT * FROM "task" WHERE "owner_id" = user_row.user_id;
    RAISE NOTICE 'First Name: %', user_row.user_id;
  END LOOP;
END $$;
