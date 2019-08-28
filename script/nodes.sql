DROP TABLE IF EXISTS `nodes`;
CREATE TABLE `nodes`(
    `node`      INTEGER PRIMARY KEY ASC,
    `name`      VARCHAR NOT NULL,
    `proto`     INT NOT NULL,
    `token`     VARCHAR NOT NULL,
    `heartbeat` TIMESTAMP NOT NULL DEFAULT 0,
    `active`    BOOLEAN NOT NULL
);

CREATE UNIQUE INDEX i_token ON nodes (`token`);
