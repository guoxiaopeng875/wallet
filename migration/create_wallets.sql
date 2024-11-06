-- Create wallets table
CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(20,4) NOT NULL DEFAULT 0.0000
);

ALTER TABLE IF EXISTS public.wallets OWNER to postgres;