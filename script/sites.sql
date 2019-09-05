DROP TABLE IF EXISTS `sites`;
CREATE TABLE `sites`(
    `site`      INTEGER PRIMARY KEY ASC,
    `url`       VARCHAR(256) NOT NULL,
    `active`    BOOLEAN NOT NULL
);

