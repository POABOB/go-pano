
-- +migrate Up
CREATE DATABASE IF NOT EXISTS `pano` DEFAULT CHARACTER SET utf8mb4;
ALTER DATABASE `pano` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down
DROP DATABASE `pano`;
