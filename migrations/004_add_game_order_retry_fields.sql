ALTER TABLE game_order
    ADD COLUMN retry_count INT UNSIGNED NULL DEFAULT 0 COMMENT '补偿重试次数' AFTER update_time,
    ADD COLUMN next_retry_time INT NULL DEFAULT 0 COMMENT '下次重试时间' AFTER retry_count,
    ADD COLUMN locked_until INT NULL DEFAULT 0 COMMENT 'Worker锁定截止时间' AFTER next_retry_time,
    ADD COLUMN last_error VARCHAR(500) NULL DEFAULT '' COMMENT '最近一次补偿错误' AFTER locked_until,
    ADD KEY idx_recovery (status, next_retry_time, locked_until);
