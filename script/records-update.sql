ALTER TABLE `records`
    ADD COLUMN `have_ocsp` BOOLEAN DEFAULT NULL;
ALTER TABLE `records`
    ADD COLUMN `ocsp_err`  VARCHAR(256) DEFAULT NULL;
ALTER TABLE `records`
    ADD COLUMN `ocsp_this_update` TIMESTAMP DEFAULT 0;
ALTER TABLE `records`
    ADD COLUMN `ocsp_next_update` TIMESTAMP DEFAULT 0;

