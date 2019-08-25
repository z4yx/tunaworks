DROP TABLE IF EXISTS `records`;
CREATE TABLE `records`(
    `http_code`      INT DEFAULT NULL,
    `response_time`  INT DEFAULT NULL,
    `site`           INT NOT NULL,
    `node`           INT NOT NULL,
    `updated`        TIMESTAMP NOT NULL,
    `ssl_err`        VARCHAR DEFAULT NULL,
    `ssl_expire`     TIMESTAMP NOT NULL
);

CREATE INDEX i_site_node_updated ON records (`site`,`node`,`updated` DESC);
