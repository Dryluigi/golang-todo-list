CREATE TABLE IF NOT EXISTS `tbl_todos` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `todo` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `is_complete` TINYINT DEFAULT 0,
    `user_id` BIGINT UNSIGNED
);