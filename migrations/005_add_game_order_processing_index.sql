ALTER TABLE game_order
    ADD KEY idx_transfer_processing (uid, agent_id, wallet_mode, status);
