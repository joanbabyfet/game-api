ALTER TABLE game_order
    ADD COLUMN wallet_mode TINYINT NULL
    COMMENT '订单创建时的钱包模式 1单一钱包 2转账钱包'
    AFTER agent_id;

-- Historical rows can only be approximated from the Agent's current mode.
UPDATE game_order AS go
INNER JOIN agent AS a ON a.id = go.agent_id
SET go.wallet_mode = a.wallet_mode
WHERE go.wallet_mode IS NULL;

UPDATE game_order
SET wallet_mode = 1
WHERE wallet_mode IS NULL;

ALTER TABLE game_order
    MODIFY COLUMN wallet_mode TINYINT NOT NULL DEFAULT 1
    COMMENT '订单创建时的钱包模式 1单一钱包 2转账钱包',
    ADD KEY idx_wallet_mode_status (wallet_mode, status);
