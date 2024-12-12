-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "status" VARCHAR(20) NOT NULL DEFAULT 'active',
  "created_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "chains" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" VARCHAR(255) NOT NULL,
  "chain_id" INT UNIQUE NOT NULL,
  "rpc_url" VARCHAR(255) NOT NULL,
  "native_currency" VARCHAR(255) NOT NULL,
  "explorer_url" VARCHAR(255),
  "status" VARCHAR(20) NOT NULL DEFAULT 'active',
  "created_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "wallets" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" UUID NOT NULL,
  "address" VARCHAR(42) NOT NULL,
  "encrypted_private_key" BYTEA NOT NULL,
  "name" VARCHAR(255),
  "status" VARCHAR(20) NOT NULL DEFAULT 'active',
  "created_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "tokens" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "chain_id" UUID NOT NULL,
  "contract_address" VARCHAR(42) NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "symbol" VARCHAR(10) NOT NULL,
  "decimals" INT NOT NULL,
  "logo_url" VARCHAR(255),
  "type" VARCHAR(20) NOT NULL DEFAULT 'ERC20',
  "status" VARCHAR(20) NOT NULL DEFAULT 'active',
  "created_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "transactions" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "chain_id" INT NOT NULL,
  "from_address" VARCHAR(42) NOT NULL,
  "to_address" VARCHAR(42) NOT NULL,
  "tx_hash" VARCHAR(66) NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE INDEX "idx_users_email" ON "users" ("email");

CREATE INDEX "idx_chains_chain_id" ON "chains" ("chain_id");

CREATE UNIQUE INDEX "unique_user_address" ON "wallets" ("user_id", "address");

CREATE INDEX "idx_wallets_user_id" ON "wallets" ("user_id");

CREATE INDEX "idx_wallets_address" ON "wallets" ("address");

CREATE UNIQUE INDEX "unique_token_per_chain" ON "tokens" ("chain_id", "contract_address");

CREATE INDEX "idx_tokens_chain_contract" ON "tokens" ("chain_id", "contract_address");

CREATE INDEX "idx_transactions_tx_hash" ON "transactions" ("tx_hash");

CREATE INDEX "idx_transactions_from_address" ON "transactions" ("from_address");

CREATE INDEX "idx_transactions_to_address" ON "transactions" ("to_address");

CREATE INDEX "idx_transactions_created_at" ON "transactions" ("created_at");


ALTER TABLE "wallets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tokens" ADD FOREIGN KEY ("chain_id") REFERENCES "chains" ("id");

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE "transactions" CASCADE;
DROP TABLE "tokens" CASCADE;
DROP TABLE "wallets" CASCADE;
DROP TABLE "chains" CASCADE;
DROP TABLE "users" CASCADE;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
