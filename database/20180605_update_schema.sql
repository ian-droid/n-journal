PRAGMA foreign_keys = 0;

CREATE TABLE sqlitestudio_temp_table AS SELECT *
                                          FROM bank;

DROP TABLE bank;

CREATE TABLE bank (
    id          INTEGER PRIMARY KEY,
    name        VARCHAR NOT NULL,
    account     VARCHAR NOT NULL,
    credit      BOOLEAN NOT NULL,
    description VARCHAR,
    active      BOOLEAN NOT NULL
                        DEFAULT (1),
    priority    BOOLEAN NOT NULL
                        DEFAULT (0),
    added       INTEGER DEFAULT (strftime('%s', 'now') ),
    updated     INTEGER
);

INSERT INTO bank (
                     id,
                     name,
                     account,
                     credit,
                     description,
                     added,
                     updated
                 )
                 SELECT id,
                        name,
                        account,
                        credit,
                        description,
                        added,
                        updated
                   FROM sqlitestudio_temp_table;

DROP TABLE sqlitestudio_temp_table;

PRAGMA foreign_keys = 1;

UPDATE bank SET active = 0 WHERE id = 2;
UPDATE bank SET priority = 1 WHERE id = 6;

UPDATE payment SET priority = 0;

PRAGMA foreign_keys = 0;

CREATE TABLE sqlitestudio_temp_table AS SELECT *
                                          FROM payment;

DROP TABLE payment;

CREATE TABLE payment (
    id          INTEGER PRIMARY KEY,
    name        VARCHAR NOT NULL,
    description VARCHAR,
    priority    BOOLEAN DEFAULT (0)
                        NOT NULL,
    added       INTEGER DEFAULT (strftime('%s', 'now') ),
    updated     INTEGER
);

INSERT INTO payment (
                        id,
                        name,
                        description,
                        priority,
                        added,
                        updated
                    )
                    SELECT id,
                           name,
                           description,
                           priority,
                           added,
                           updated
                      FROM sqlitestudio_temp_table;

DROP TABLE sqlitestudio_temp_table;

PRAGMA foreign_keys = 1;

UPDATE payment SET priority = 1 WHERE id = 3;

UPDATE currency SET prefix = "CN￥" WHERE id = 1;

UPDATE bank SET name = "Balance", WHERE id = 0;
UPDATE bank SET name = "ICBC借记", WHERE id = 1;
