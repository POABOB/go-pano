
-- +migrate Up
CREATE TABLE IF NOT EXISTS `Clinic` (
  `clinic_id` int PRIMARY KEY AUTO_INCREMENT COMMENT '診所ID',
  `name` varchar(128) NOT NULL COMMENT '診所名稱',
  `start_at` date NOT NULL COMMENT '開始日',
  `end_at` date NOT NULL DEFAULT "9999-12-31" COMMENT '結束日',
  `quota_per_month` int NOT NULL DEFAULT 200 COMMENT '每月使用額度',
  `token` varchar(1024) NOT NULL COMMENT 'Token'
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `Record` (
  `record_id` int PRIMARY KEY AUTO_INCREMENT COMMENT '紀錄ID',
  `clinic_id` int DEFAULT 0 COMMENT '診所ID',
  `predict_id` int DEFAULT 0 COMMENT '預測ID',
  `score` int NOT NULL DEFAULT 80 COMMENT '準確度',
  `comment` varchar(1024) NOT NULL DEFAULT "" COMMENT '準確度評論'
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `Users` (
  `user_id` int PRIMARY KEY AUTO_INCREMENT COMMENT '使用者ID',
  `name` varchar(64) DEFAULT "" COMMENT '帳號名稱',
  `account` varchar(64) DEFAULT "" COMMENT '帳號',
  `password` varchar(1024) DEFAULT "" COMMENT '密碼',
  `roles_string` varchar(1024) DEFAULT "[]" COMMENT '權限',
  `status` int DEFAULT 1 COMMENT '啟用狀態'
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `Predict` CHANGE `predict_id` `predict_id` int NOT NULL AUTO_INCREMENT COMMENT '辨識ID';
ALTER TABLE `Predict` CHANGE `clinic_id` `clinic_id` int not null DEFAULT 0 COMMENT '診所ID';
ALTER TABLE `Predict` CHANGE `filename` `filename` varchar(128) not null DEFAULT '' COMMENT '檔案名稱';
ALTER TABLE `Predict` CHANGE `predict` `predict` varchar(1024) not null DEFAULT '' COMMENT '辨識';
ALTER TABLE `Predict` CHANGE `created_at` `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '新增時間';
ALTER TABLE `Predict` CHANGE `updated_at` `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新時間';


CREATE INDEX `Clinic_index_0` ON `Clinic` (`start_at`);

CREATE INDEX `Clinic_index_1` ON `Clinic` (`end_at`);

CREATE INDEX `Clinic_index_2` ON `Clinic` (`token`);

CREATE INDEX `Predict_index_3` ON `Predict` (`clinic_id`);

CREATE INDEX `Predict_index_4` ON `Predict` (`filename`);

CREATE INDEX `Record_index_5` ON `Record` (`clinic_id`);

CREATE INDEX `Record_index_6` ON `Record` (`predict_id`);

CREATE INDEX `User_index_7` ON `Users` (`user_id`);

CREATE INDEX `User_index_8` ON `Users` (`account`);

CREATE INDEX `User_index_9` ON `Users` (`password`);

CREATE INDEX `User_index_10` ON `Users` (`status`);

ALTER TABLE `Predict` ADD CONSTRAINT `fk_clinic_id_1` FOREIGN KEY (`clinic_id`) REFERENCES `Clinic` (`clinic_id`);

ALTER TABLE `Record` ADD CONSTRAINT `fk_clinic_id_2` FOREIGN KEY (`clinic_id`) REFERENCES `Clinic` (`clinic_id`);

ALTER TABLE `Record` ADD CONSTRAINT `fk_predict_id_1` FOREIGN KEY (`predict_id`) REFERENCES `Predict` (`predict_id`);



-- +migrate Down
ALTER TABLE `Predict` DROP FOREIGN KEY `fk_clinic_id_1`;
ALTER TABLE `Record` DROP FOREIGN KEY `fk_clinic_id_2`;
ALTER TABLE `Record` DROP FOREIGN KEY `fk_predict_id_1`;
DROP TABLE IF EXISTS `Clinic`;
DROP TABLE IF EXISTS `Record`;
DROP TABLE IF EXISTS `Users`;

-- TEST PREDICT
-- REAL TEST USER
-- REAL TEST PREDICT