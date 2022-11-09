
-- +migrate Up
-- 新增Predict資料表
CREATE TABLE IF NOT EXISTS `Predict` (
  `id` int NOT NULL AUTO_INCREMENT,
  `clinic_id` int NOT NULL DEFAULT 0,
  `filename` varchar(100) DEFAULT '',
  `predict` varchar(1024) DEFAULT '',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `Predict`;
