
-- +migrate Up
-- 新增Predict資料表
CREATE TABLE IF NOT EXISTS `Predict` (
  `predict_id` int NOT NULL AUTO_INCREMENT,
  `clinic_id` int NOT NULL DEFAULT 0,
  `dir` varchar(64) NOT NULL DEFAULT '',
  `filename` varchar(100) DEFAULT '',
  `predict_string` varchar(1024) DEFAULT '',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`predict_id`)
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `Predict`;
