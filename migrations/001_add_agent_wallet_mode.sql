ALTER TABLE agent
    ADD COLUMN wallet_mode TINYINT NULL DEFAULT 1 COMMENT '钱包模式 1单一钱包 2转账钱包'
    AFTER callback_url,
    ADD KEY idx_wallet_mode (wallet_mode);
