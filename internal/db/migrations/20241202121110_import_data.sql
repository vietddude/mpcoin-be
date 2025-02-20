-- +goose Up
INSERT INTO chains (id, chain_id, name, native_currency, rpc_url, explorer_url, created_at, updated_at) VALUES ('6397f2b3-b503-458f-8df2-df1f6719266a', '11155111', 'Sepolia','ETH', 'https://ethereum-sepolia-rpc.publicnode.com', 'https://sepolia.etherscan.io', NOW(), NOW());
INSERT INTO tokens 
    (chain_id, contract_address, name, symbol, decimals, logo_url, type, created_at, updated_at) 
VALUES 
    ('6397f2b3-b503-458f-8df2-df1f6719266a', '0x1111111111111111111111111111111111111111', 'Ethereum', 'ETH', 18, 'https://arbiscan.io/token/images/ether.svg', 'NATIVE', NOW(), NOW()),
    ('6397f2b3-b503-458f-8df2-df1f6719266a', '0xFB122130C4d28860dbC050A8e024A71a558eB0C1', 'USDT', 'USDT', 18, 'https://s2.coinmarketcap.com/static/img/coins/64x64/825.png', 'ERC20', NOW(), NOW());
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
