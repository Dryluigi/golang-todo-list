ALTER TABLE `tbl_todos`
ADD CONSTRAINT `fk_todos_users`
FOREIGN KEY (`user_id`) REFERENCES `tbl_users` (`id`);