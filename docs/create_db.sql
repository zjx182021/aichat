create database ai_chat default charset utf8mb4;
use ai_chat;


CREATE TABLE `chat_records` (
 `id` bigint NOT NULL AUTO_INCREMENT,
 `user_msg` text,
 `user_msg_tokens` int NOT NULL DEFAULT '0',
 `user_msg_keywords` varchar(1024) NOT NULL DEFAULT '',
 `ai_msg` text,
 `ai_msg_tokens` int NOT NULL DEFAULT '0',
 `req_tokens` int NOT NULL DEFAULT '0',
 `create_at` bigint NOT NULL DEFAULT '0',
 PRIMARY KEY (`id`),
 KEY `index_create_at` (`create_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci