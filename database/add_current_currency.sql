PRAGMA foreign_keys = 0;

CREATE TABLE sqlitestudio_temp_table AS SELECT *
                                          FROM currency;

DROP TABLE currency;

CREATE TABLE currency (
    id      INTEGER     PRIMARY KEY
                        NOT NULL,
    name    VARCHAR     NOT NULL,
    code    CHAR (3)    NOT NULL,
    prefix  VARCHAR (5) NOT NULL,
    current BOOLEAN     NOT NULL
                        DEFAULT (0)
);

INSERT INTO currency (
                         id,
                         name,
                         code,
                         prefix
                     )
                     SELECT id,
                            name,
                            code,
                            prefix
                       FROM sqlitestudio_temp_table;

DROP TABLE sqlitestudio_temp_table;

PRAGMA foreign_keys = 1;
