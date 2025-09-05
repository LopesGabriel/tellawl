DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_wallet_id;
DROP INDEX IF EXISTS idx_wallet_users_user_id;
DROP INDEX IF EXISTS idx_wallets_creator_id;

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallet_users;
DROP TABLE IF EXISTS wallets;