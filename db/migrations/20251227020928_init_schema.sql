-- +goose Up
-- +goose StatementBegin

-- 1. Config: One row per guild.
CREATE TABLE guild_config (
    guild_id TEXT PRIMARY KEY,
    starting_cash BIGINT NOT NULL DEFAULT 100000, -- e.g. 1000.00 cents
    base_price BIGINT NOT NULL DEFAULT 1000,      -- 10.00 cents
    slope BIGINT NOT NULL DEFAULT 100,            -- 1.00 cent increase per share
    trade_fee_bps INT NOT NULL DEFAULT 200,       -- 2%
    subject_fee_bps INT NOT NULL DEFAULT 100,     -- 1% (half of trade fee)
    min_trade_qty BIGINT NOT NULL DEFAULT 1,
    sunset_days INT NOT NULL DEFAULT 7,
    trading_enabled BOOLEAN NOT NULL DEFAULT true,
    -- Immutable by convention (enforced in App logic):
    currency_unit TEXT NOT NULL DEFAULT 'cents' CHECK (currency_unit = 'cents')
);

-- 2. Members: The registry of users in the economy.
CREATE TABLE guild_members (
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    opted_in BOOLEAN NOT NULL DEFAULT false,
    opted_in_at TIMESTAMPTZ,
    opted_out_at TIMESTAMPTZ,
    PRIMARY KEY (guild_id, user_id)
);

-- 3. Wallets: Holds user cash.
CREATE TABLE wallets (
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    cash BIGINT NOT NULL DEFAULT 0 CHECK (cash >= 0),
    PRIMARY KEY (guild_id, user_id),
    FOREIGN KEY (guild_id, user_id) REFERENCES guild_members(guild_id, user_id) ON DELETE RESTRICT
);

-- 4. Markets: The bonding curve state for a specific subject.
CREATE TABLE markets (
    guild_id TEXT NOT NULL,
    subject_user_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('ACTIVE', 'SUNSETTING', 'CLOSED')),
    shares_outstanding BIGINT NOT NULL DEFAULT 0 CHECK (shares_outstanding >= 0),
    reserve_balance BIGINT NOT NULL DEFAULT 0 CHECK (reserve_balance >= 0),
    last_price BIGINT NOT NULL,
    sunset_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (guild_id, subject_user_id),
    FOREIGN KEY (guild_id, subject_user_id) REFERENCES guild_members(guild_id, user_id) ON DELETE RESTRICT
);

-- 5. Treasury: Collects system fees.
CREATE TABLE treasury (
    guild_id TEXT PRIMARY KEY,
    system_fees BIGINT NOT NULL DEFAULT 0 CHECK (system_fees >= 0)
);

-- 6. Holdings: Who owns what.
CREATE TABLE holdings (
    guild_id TEXT NOT NULL,
    owner_user_id TEXT NOT NULL,
    subject_user_id TEXT NOT NULL,
    shares BIGINT NOT NULL DEFAULT 0 CHECK (shares >= 0),
    avg_cost BIGINT, -- Nullable, purely informational
    PRIMARY KEY (guild_id, owner_user_id, subject_user_id),
    FOREIGN KEY (guild_id, owner_user_id) REFERENCES guild_members(guild_id, user_id) ON DELETE RESTRICT,
    FOREIGN KEY (guild_id, subject_user_id) REFERENCES markets(guild_id, subject_user_id) ON DELETE RESTRICT
);

-- 7. Trades: The immutable ledger of all actions.
CREATE TABLE trades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    guild_id TEXT NOT NULL,
    trader_user_id TEXT NOT NULL,
    subject_user_id TEXT NOT NULL,
    side TEXT NOT NULL CHECK (side IN ('BUY', 'SELL')),
    shares BIGINT NOT NULL,
    
    -- Financial breakdown
    gross_amount BIGINT NOT NULL,
    fee_amount BIGINT NOT NULL,
    subject_fee_amount BIGINT NOT NULL,
    system_fee_amount BIGINT NOT NULL,
    net_amount BIGINT NOT NULL,
    
    -- State snapshots for audit
    price_before BIGINT NOT NULL,
    price_after BIGINT NOT NULL,
    shares_outstanding_after BIGINT NOT NULL,
    reserve_balance_after BIGINT NOT NULL,
    
    -- Idempotency & Meta
    idempotency_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    UNIQUE (idempotency_key),
    FOREIGN KEY (guild_id, trader_user_id) REFERENCES guild_members(guild_id, user_id) ON DELETE RESTRICT,
    FOREIGN KEY (guild_id, subject_user_id) REFERENCES markets(guild_id, subject_user_id) ON DELETE RESTRICT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trades;
DROP TABLE IF EXISTS holdings;
DROP TABLE IF EXISTS treasury;
DROP TABLE IF EXISTS markets;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS guild_members;
DROP TABLE IF EXISTS guild_config;
-- +goose StatementEnd