-- +goose Up
-- +goose StatementBegin

CREATE INDEX idx_markets_price ON markets(guild_id, last_price DESC) WHERE status = 'ACTIVE';
CREATE INDEX idx_holdings_owner ON holdings(guild_id, owner_user_id);
CREATE INDEX idx_trades_history ON trades(guild_id, subject_user_id, created_at DESC);
CREATE INDEX idx_trades_trader ON trades(guild_id, trader_user_id, created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_trades_trader;
DROP INDEX IF EXISTS idx_trades_history;
DROP INDEX IF EXISTS idx_holdings_owner;
DROP INDEX IF EXISTS idx_markets_price;

-- +goose StatementEnd
