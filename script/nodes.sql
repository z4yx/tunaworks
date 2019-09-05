DROP TABLE IF EXISTS `nodes`;
CREATE TABLE `nodes`(
    `node`      INTEGER PRIMARY KEY,
    `name`      VARCHAR(64) NOT NULL,
    `proto`     INT NOT NULL,
    `token`     VARCHAR(64) NOT NULL,
    `heartbeat` TIMESTAMP NOT NULL DEFAULT 0,
    `active`    BOOLEAN NOT NULL
);

CREATE UNIQUE INDEX i_token ON nodes (`token`);
