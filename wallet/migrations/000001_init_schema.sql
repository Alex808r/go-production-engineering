-- Enable UUID extension if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table for wallets
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY,
    balance NUMERIC(18, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for fast lookup by wallet id
CREATE INDEX IF NOT EXISTS idx_wallets_id ON wallets(id);

-- Optional audit log for wallet transactions
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    operation_type VARCHAR(20) NOT NULL,
    amount NUMERIC(18, 2) NOT NULL,
    balance_after NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallet_transactions_wallet_id ON wallet_transactions(wallet_id);

-- Default demo wallet for instant testing
INSERT INTO wallets (id, balance)
VALUES ('11111111-1111-1111-1111-111111111111', 10000.00)
ON CONFLICT (id) DO NOTHING;
