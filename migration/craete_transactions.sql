-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    method VARCHAR(10) NOT NULL,
    tx_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    amount DECIMAL(20,2) NOT NULL,
    from_wallet_id INTEGER,
    to_wallet_id INTEGER
);

ALTER TABLE IF EXISTS public.transactions OWNER to postgres;