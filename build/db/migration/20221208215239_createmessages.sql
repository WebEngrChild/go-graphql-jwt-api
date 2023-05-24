-- +goose Up
CREATE TABLE messages (
   id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
   Text VARCHAR(255) DEFAULT NULL,
   user_id INT UNSIGNED NOT NULL,
   created_at TIMESTAMP DEFAULT NULL,
   updated_at TIMESTAMP DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE INDEX message_id on messages (id);

INSERT INTO messages (Text, user_id, created_at)
VALUES ('message1,message1,message1,message1,message1,message1,message1,message1,message1,message1,message1,', '1', DATE_SUB(NOW(), INTERVAL FLOOR(RAND()*365) DAY)),
       ('message2,message2,message2,message2,message2,message2,message2,message2,message2,message2,message2,', '2', DATE_SUB(NOW(), INTERVAL FLOOR(RAND()*365) DAY)),
       ('message3,message3,message3,message3,message3,message3,message3,message3,message3,message3,message3,', '3', DATE_SUB(NOW(), INTERVAL FLOOR(RAND()*365) DAY));

-- +goose Down
DROP INDEX message_id;
DROP TABLE messages;