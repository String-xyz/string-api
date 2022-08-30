-- +goose Up
CREATE TABLE transaction (
    id int NOT NULL,
    created date NOT NULL,
    entityType int, --foreign key
    chainID int,
    userAddress text,
    contractAddress text,
    contractParameters text, -- can't take slice here
    contractABI text, -- can't take slice here
    contractFunction text,
    txValue text,
    gasLimit text,
    cardToken text,
    status int, --foreign key (TransactionStatus)
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;