--
-- Database :  `rest_api`
--
CREATE DATABASE IF NOT EXISTS `rest_api` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `rest_api`;

-- Drop child table first to avoid FK constraint errors
DROP TABLE IF EXISTS `user_task`;
DROP TABLE IF EXISTS `task`;
DROP TABLE IF EXISTS `user`;

CREATE TABLE IF NOT EXISTS `task` (
  `task_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(128) NOT NULL,
  `description` text NOT NULL,
  `creation_date` datetime NOT NULL,
  `status` tinyint(3) NOT NULL DEFAULT '1',
  PRIMARY KEY (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS `user` (
  `user_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `email` varchar(320) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS `user_task` (
  `user_id` int(11) unsigned NOT NULL,
  `task_id` int(11) unsigned NOT NULL,
  UNIQUE KEY `user_task` (`user_id`,`task_id`),
  CONSTRAINT `fk_ut_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ut_task` FOREIGN KEY (`task_id`) REFERENCES `task` (`task_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `rest_api`.`user` (`user_id`, `name`, `email`) VALUES ('1', 'G4', 'gectou4@gmail.com');

INSERT INTO `rest_api`.`task` (`task_id`, `title`, `description`, `creation_date`, `status`) VALUES ('1', 'Faire le cafée', 'Aller à la cafetière
mettre la tasse
Mettre le café
Mettre de l''eau si besoin
Appuyer sur le bouton "Café tiptop"
Attendre que l''eau ne goutte plus
Prendre la tasse', '2015-11-04 03:01:08', '2'), ('2', 'Boire le café', 'Prendre la tasse de café de la Task 1
Savourer', '2015-11-04 03:01:08', '1');

INSERT INTO `rest_api`.`user_task` (`user_id`, `task_id`) VALUES ('1', '1'), ('1', '2');
