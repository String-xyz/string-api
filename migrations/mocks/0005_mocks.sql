-------------------------------------------------------------------------
-- +goose Up
-------------------------------------------------------------------------
-- STRING_USER ----------------------------------------------------------
INSERT INTO string_user (id, created_at, updated_at, type, status, tags, first_name, middle_name, last_name)
VALUES ('0e837b73-55cf-43ff-9b1e-0d8258eec978', '2022-10-19 00:17:01.837572+00', '2022-10-19 00:17:01.837572+00', 'Developer', 'Developing', '{}', 'Deve', '', 'Loper');

-------------------------------------------------------------------------
-- DEVICE ---------------------------------------------------------------
INSERT INTO device (id, created_at, updated_at, last_used_at, validated_at, type, description, fingerprint, user_id)
VALUES ('073f5a88-9223-4554-a7ce-11d358123a21', '2022-10-19 00:23:10.405595+00', '2022-10-19 00:23:10.405595+00', '2022-10-19 00:17:01.837572+00', '2022-10-19 00:17:01.837572+00', '', 'Developer Laptop', '', '0e837b73-55cf-43ff-9b1e-0d8258eec978');

-------------------------------------------------------------------------
-- INSTRUMENT -----------------------------------------------------------
INSERT INTO instrument (id, created_at, updated_at, type, status, tags, network, public_key, last_4, user_id)
VALUES ('13438963-f5e7-47c4-a790-ebca3e3bf915', '2022-10-19 00:53:36.538289+00', '2022-10-19 00:53:36.538289+00', 'Credit Card', 'Ephemeral', '{}', 'Mastercard', '', '4242', '0e837b73-55cf-43ff-9b1e-0d8258eec978'), 
('ab6a2d66-ad4c-43f4-adf9-c0cd3282492c', '2022-10-19 00:55:26.166175+00', '2022-10-19 00:55:26.166175+00', 'Crypto Wallet', 'Ephemeral', '{}', 'Ethereum', '0x44A4b9E2A69d86BA382a511f845CbF2E31286770', '', '0e837b73-55cf-43ff-9b1e-0d8258eec978');

-------------------------------------------------------------------------
-- NETWORK --------------------------------------------------------------
INSERT INTO network (id, created_at, updated_at, name, network_id, chain_id, gas_token_id, gas_oracle, rpc_url)
VALUES ('ea34e526-ec6e-4f2b-89b4-acc08db80d63', '2022-10-14 20:18:09.555645+00', '2022-10-14 20:18:09.555645+00', 'Fuji Testnet', '43113', '43113', '19611d0e-a42f-4cee-a35a-b34eb5c08a7f', 'avax', 'https://api.avax-test.network/ext/bc/C/rpc'), 
('b21d6cd6-5d8a-49a6-bac6-e6323316dc01', '2022-10-14 20:41:39.962327+00', '2022-10-14 20:41:39.962327+00', 'Goerli Testnet', '5', '5', '3ef72571-c2e1-4ca3-991c-0df17cef7535', 'eth', 'https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161'),
('6cea71b3-b287-4680-ad9d-e631d0bc84ba', '2022-10-14 20:41:39.962327+00', '2022-10-14 20:41:39.962327+00', 'Polygon Mainnet', '137', '137', 'c06986d8-cc2c-4cdc-9728-16a45698b3e7', 'poly', 'https://rpc-mainnet.matic.quiknode.pro'),
('cd42c066-554c-42ad-994b-48fed371931c', '2022-10-14 20:41:39.962327+00', '2022-10-14 20:41:39.962327+00', 'Avalanche Mainnet', '43114', '43114', '19611d0e-a42f-4cee-a35a-b34eb5c08a7f', 'avax', 'https://api.avax.network/ext/bc/C/rpc'),
('491d46e2-18e0-45ec-8209-faf0ec5d278c', '2022-10-14 20:41:39.962327+00', '2022-10-14 20:41:39.962327+00', 'Mumbai Testnet', '80001', '80001', 'c06986d8-cc2c-4cdc-9728-16a45698b3e7', 'poly', 'https://matic-mumbai.chainstacklabs.com'),
('60a02818-4e7d-4b84-b673-e2376fdbfbf9', '2022-10-14 20:41:39.962327+00', '2022-10-30 01:07:37.237054+00', 'Ethereum Mainnet', '1', '1', '3ef72571-c2e1-4ca3-991c-0df17cef7535', 'eth', 'https://rpc.ankr.com/eth');

-------------------------------------------------------------------------
-- PLATFORM -------------------------------------------------------------
INSERT INTO platform (id, created_at, updated_at, type, status, name, api_key, authentication)
VALUES ('54a7e062-4cec-44f3-9d89-99498d0eb6ef', '2022-10-19 00:37:16.965408+00', '2022-10-19 00:37:16.965408+00', 'Game', 'Verified', 'Nintendo', 'developer', 'email');

-------------------------------------------------------------------------
-- ASSET -------------------------------------------------------------
INSERT INTO asset (id, created_at, updated_at, name, description, decimals, is_crypto, network_id, value_oracle)
VALUES ('19611d0e-a42f-4cee-a35a-b34eb5c08a7f', '2022-10-14 20:17:06.460812+00', '2022-10-15 02:41:02.270712+00', 'AVAX', 'Avalanche', 18, TRUE, 'cd42c066-554c-42ad-994b-48fed371931c', 'avalanche-2'),
('c06986d8-cc2c-4cdc-9728-16a45698b3e7', '2022-10-14 20:17:06.460812+00', '2022-10-15 02:41:02.270712+00', 'MATIC', 'Matic', 18, TRUE, '6cea71b3-b287-4680-ad9d-e631d0bc84ba', 'matic-network'),
('3ef72571-c2e1-4ca3-991c-0df17cef7535', '2022-10-14 20:17:06.460812+00', '2022-10-15 02:41:02.270712+00', 'ETH', 'Ethereum', 18, TRUE, '60a02818-4e7d-4b84-b673-e2376fdbfbf9', 'ethereum'),
('bc376c3a-6481-49d0-83ef-34ba80937ba8', '2022-10-18 03:59:05.042924+00', '2022-10-18 03:59:05.042924+00', 'USD', 'United States Dollar', 6, FALSE, null, '');


-------------------------------------------------------------------------
-- +goose Down

-- Can't delete rows due to foreign key constraint