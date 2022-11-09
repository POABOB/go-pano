
-- +migrate Up
-- 新增pano DB
CREATE DATABASE IF NOT EXISTS `pano` DEFAULT CHARACTER SET utf8mb4;

-- +migrate Down
-- DROP DATABASE `pano`;
