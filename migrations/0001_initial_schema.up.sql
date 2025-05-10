CREATE TABLE `transactions` (
  `id` char(36) NOT NULL,
  `sender` longtext,
  `receiver` longtext,
  `amount` double DEFAULT NULL,
  `is_mined` tinyint(1) DEFAULT '0',
  `block_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `blocks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `previous_hash` longtext,
  `hash` longtext,
  `nonce` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_blocks_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;