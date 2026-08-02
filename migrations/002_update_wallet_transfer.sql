ALTER TABLE wallet_transfer
    ADD COLUMN request_id VARCHAR(64) NULL DEFAULT '' COMMENT 'Provider请求ID' AFTER order_no,
    ADD COLUMN third_order_no VARCHAR(64) NULL DEFAULT '' COMMENT '第三方订单号' AFTER request_id,
    ADD COLUMN currency VARCHAR(10) NULL DEFAULT '' COMMENT '币种' AFTER amount,
    ADD COLUMN balance_before BIGINT NULL DEFAULT 0 COMMENT '转账前余额' AFTER currency,
    ADD COLUMN balance_after BIGINT NULL DEFAULT 0 COMMENT '转账后余额' AFTER balance_before,
    ADD COLUMN error_code INT NULL DEFAULT 0 COMMENT '失败错误码' AFTER status,
    ADD COLUMN error_message VARCHAR(255) NULL DEFAULT '' COMMENT '失败原因' AFTER error_code,
    ADD COLUMN finish_time INT NULL DEFAULT 0 COMMENT '完成时间' AFTER update_time;

-- Existing rows need non-empty unique values before the unique indexes are added.
UPDATE wallet_transfer
SET request_id = CONCAT('legacy-', order_no),
    third_order_no = CONCAT('legacy-', order_no)
WHERE request_id = '' OR third_order_no = '';

ALTER TABLE wallet_transfer
    ADD UNIQUE KEY uk_agent_third_order (agent_id, third_order_no),
    ADD UNIQUE KEY uk_agent_request (agent_id, request_id);
