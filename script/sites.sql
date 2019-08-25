DROP TABLE IF EXISTS `sites`;
CREATE TABLE `sites`(
    `site`      INTEGER PRIMARY KEY ASC,
    `url`       VARCHAR NOT NULL,
    `active`    BOOLEAN NOT NULL
);

