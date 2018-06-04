BEGIN;
UPDATE diary SET highlighted = 0 where highlighted is NULL;
UPDATE diary SET highlighted = 1 WHERE highlighted= "TRUE" or highlighted;
UPDATE diary SET highlighted = 0 WHERE not highlighted;

PRAGMA foreign_keys = 0;

CREATE TABLE sqlitestudio_temp_table AS SELECT *
                                          FROM diary;

DROP TABLE diary;

CREATE TABLE diary (
    date        DATE    PRIMARY KEY NOT NULL,
    content     VARCHAR NOT NULL,
    highlighted BOOLEAN DEFAULT 0
                        NOT NULL,
    added       INTEGER DEFAULT (strftime('%s', 'now') ),
    updated     INTEGER
);

INSERT INTO diary (
                      date,
                      content,
                      highlighted,
                      added,
                      updated
                  )
                  SELECT date,
                         content,
                         highlighted,
                         added,
                         updated
                    FROM sqlitestudio_temp_table;

DROP TABLE sqlitestudio_temp_table;

PRAGMA foreign_keys = 1;

Commit;
