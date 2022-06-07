CREATE TABLE `file_upload`
(
    `id`           INT(11)       NOT NULL AUTO_INCREMENT,
    `file_path`    VARCHAR(500)  NULL DEFAULT '' COLLATE 'utf8_unicode_ci',
    `state`        TINYINT(1)    NULL DEFAULT '1' COMMENT '1: processing 2: success 3: failed',
    `description`  VARCHAR(1024) NULL DEFAULT '' COLLATE 'utf8_unicode_ci',
    `created_time` INT(11)       NULL DEFAULT '0',
    `updated_time` INT(11)       NULL DEFAULT '0',
    `created_at`   TIMESTAMP     NULL DEFAULT current_timestamp(),
    `updated_at`   TIMESTAMP     NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `state` (`state`) USING BTREE
)
    COLLATE = 'utf8_vietnamese_ci'
    ENGINE = InnoDB
    AUTO_INCREMENT = 17
;
