WITH "new_users" AS
(
  INSERT INTO "user"
              ("first_name", "middle_name", "last_name",   "surname",   "picture_url",                                     "email",                    "password")
       VALUES ('Briano',     'Geoff',       'Cakebread',   'Frie',      'http://dummyimage.com/235x100.png/5fa2dd/ffffff', 'gfrie0@addtoany.com',      '5y2S0QM'),
              ('Harmonie',   'Merla',       'Saberton',    'Shoulders', 'http://dummyimage.com/163x100.png/cc0000/ffffff', 'mshoulders1@shop-pro.jp',  'Rbdr8l'),
              ('Roscoe',     'Merry',       'Sibyllina',   'Dixson',    'http://dummyimage.com/145x100.png/cc0000/ffffff', 'mdixson2@typepad.com',     'Xk1sjdPgey'),
              ('Nedda',      'Kristin',     'Lewin',       'Crispin',   'http://dummyimage.com/207x100.png/5fa2dd/ffffff', 'kcrispin3@alexa.com',      'FdVWjWbfN0'),
              ('Noelyn',     'Muriel',      'De Few',      'Fewkes',    'http://dummyimage.com/165x100.png/5fa2dd/ffffff', 'mfewkes4@photobucket.com', 'M8mgK.kDBp'),
              ('Austin',     'Skyler',      'Kitchenside', 'Masson',    'http://dummyimage.com/153x100.png/dddddd/000000', 'smasson5@blogspot.com',    'nFHMVqFlVJ'),
              ('Shirlene',   'Illa',        'Staynes',     'MacAless',  'http://dummyimage.com/176x100.png/cc0000/ffffff', 'imacaless6@ameblo.jp',     'vFBFYZjg7d'),
              ('Sherrie',    'Hamnet',      'Prestedge',   'Fackney',   'http://dummyimage.com/214x100.png/5fa2dd/ffffff', 'hfackney7@patch.com',      'PmfX'),
              ('Jared',      'Catlaina',    'McFarlane',   'Craighill', 'http://dummyimage.com/212x100.png/ff4444/ffffff', 'ccraighill8@blogs.com',    'yW9sAf'),
              ('Eba',        'Raynard',     'Yakovl',      'Gurnett',   'http://dummyimage.com/182x100.png/dddddd/000000', 'rgurnett9@guardian.co.uk', '7yRXvsngK7c')
    RETURNING "user_id", "first_name", "last_name"
)
INSERT INTO "list" ("owner_id", "name", "description")
     SELECT "user_id",
            'today',
            "first_name" || ' ' || "last_name" || '''s today list'
       FROM "new_users"
  UNION ALL
     SELECT "user_id",
            'tomorrow',
            "first_name" || ' ' || "last_name" || '''s tomorrow list'
       FROM "new_users";
