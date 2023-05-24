-- +goose Up
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    name VARCHAR(255) DEFAULT NULL,
    email VARCHAR(255) DEFAULT NULL,
    password VARCHAR(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE INDEX user_id on users (id);

SET @encryption_key = 'passwordpassword';
INSERT INTO users (name, email, password)
VALUES
    ('user1', HEX(AES_ENCRYPT('test1@example.com', @encryption_key)), HEX(AES_ENCRYPT('password1', @encryption_key))),
    ('user2', HEX(AES_ENCRYPT('test2@example.com', @encryption_key)), HEX(AES_ENCRYPT('password2', @encryption_key))),
    ('user3', HEX(AES_ENCRYPT('test3@example.com', @encryption_key)), HEX(AES_ENCRYPT('password3', @encryption_key)));

-- +goose Down
DROP INDEX user_id;
DROP TABLE users;
