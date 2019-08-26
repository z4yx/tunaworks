DROP TABLE IF EXISTS `nodes`;
CREATE TABLE `nodes`(
    `node`      INTEGER PRIMARY KEY ASC,
    `name`      VARCHAR NOT NULL,
    `proto`     INT NOT NULL,
    `token`     VARCHAR NOT NULL,
    `active`    BOOLEAN NOT NULL
);

