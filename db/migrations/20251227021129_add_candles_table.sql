-- +goose Up
-- +goose StatementBegin

CREATE TABLE candles_1h (
    guild_id TEXT NOT NULL,
    subject_user_id TEXT NOT NULL,
    bucket_start TIMESTAMPTZ NOT NULL,
    open BIGINT NOT NULL,
    high BIGINT NOT NULL,
    low BIGINT NOT NULL,
    close BIGINT NOT NULL,
    volume BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (guild_id, subject_user_id, bucket_start),
    FOREIGN KEY (guild_id, subject_user_id) REFERENCES markets(guild_id, subject_user_id) ON DELETE RESTRICT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS candles_1h;

-- +goose StatementEnd
