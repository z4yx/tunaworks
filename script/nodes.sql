DROP TABLE IF EXISTS `nodes`;
CREATE TABLE `nodes`(
    `node`      INTEGER PRIMARY KEY ASC,
    `name`      VARCHAR NOT NULL,
    `active`    BOOLEAN NOT NULL
);

