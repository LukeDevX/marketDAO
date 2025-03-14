-- 1. 用户表：存储用户基本信息（不直接存储钱包地址）
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,  -- 用户名
    email VARCHAR(100) UNIQUE,  -- 邮箱
    password_hash VARCHAR(255) NOT NULL,  -- 密码哈希
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


-- 2. 用户钱包表：存储不同链的钱包地址
CREATE TABLE user_wallets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,  -- 关联用户
    chain ENUM('EVM', 'TON', 'Solana') NOT NULL,  -- 所属区块链
    wallet_address VARCHAR(255) NOT NULL,  -- 钱包地址
    UNIQUE (user_id, chain)  -- 确保同一个用户在同一条链上只有一个钱包
     -- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. NFT 资产表：每个NFT归属一个钱包，而不是用户
CREATE TABLE nfts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    contract_address VARCHAR(255) NOT NULL,  -- 合约地址
    owner_wallet_address VARCHAR(255), -- 拥有者钱包地址
    chain ENUM('EVM', 'TON', 'Solana') NOT NULL,  -- 所属区块链
    token_id VARCHAR(255) NOT NULL,  -- NFT唯一ID
    metadata_uri VARCHAR(500),  -- 元数据地址
    status ENUM('listed', 'unlisted', 'sold') NOT NULL DEFAULT 'unlisted'  -- NFT状态
    -- FOREIGN KEY (wallet_id) REFERENCES user_wallets(id) ON DELETE CASCADE
);


-- 4. NFT 交易记录表：记录交易时的钱包
CREATE TABLE nft_transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    nft_id BIGINT NOT NULL,  -- 交易的NFT
    buyer_wallet_address VARCHAR(255) NOT NULL,  -- 买家钱包ID
    seller_wallet_address VARCHAR(255) NOT NULL,  -- 卖家钱包ID
    price DECIMAL(30, 10) NOT NULL,  -- 交易金额
    tx_hash VARCHAR(255) UNIQUE NOT NULL,  -- 链上交易哈希
    chain ENUM('EVM', 'TON', 'Solana') NOT NULL,  -- 交易所属区块链
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    -- FOREIGN KEY (nft_id) REFERENCES nfts(id) ON DELETE CASCADE,
    -- FOREIGN KEY (buyer_wallet_id) REFERENCES user_wallets(id),
    -- FOREIGN KEY (seller_wallet_id) REFERENCES user_wallets(id)
);


-- 5. swap 交易记录表
CREATE TABLE swap_transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    wallet_id BIGINT NOT NULL,
    from_asset_address VARCHAR(255) NOT NULL,
    to_asset_address VARCHAR(255) NOT NULL,
    from_amount DECIMAL(18, 8) NOT NULL,
    to_amount DECIMAL(18, 8) NOT NULL,
    chain ENUM('EVM', 'TON', 'Solana') NOT NULL,  -- 所属区块链
    transaction_hash VARCHAR(255) UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    -- FOREIGN KEY (wallet_id) REFERENCES wallets(id)
);


-- 6  系统操作日志
CREATE TABLE system_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    action_type VARCHAR(255) NOT NULL,
    action_details TEXT,
    chain_type ENUM('EVM', 'TON', 'Solana'),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    -- FOREIGN KEY (user_id) REFERENCES users(id)
);
