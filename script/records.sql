DROP TABLE IF EXISTS `records`;
CREATE TABLE `records`(
    `http_code`      INT DEFAULT NULL,
    `response_time`  INT DEFAULT NULL,
    `site`           INT NOT NULL,
    `node`           INT NOT NULL,
    `protocol`       INT NOT NULL,
    `updated`        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `ssl_err`        VARCHAR(256) DEFAULT NULL,
    `ssl_expire`     TIMESTAMP NOT NULL
);

CREATE INDEX i_s_n_p_u ON records (`site`,`node`,`protocol`,`updated` DESC);
