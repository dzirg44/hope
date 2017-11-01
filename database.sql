CREATE  DATABASE gophr;

CREATE TABLE `images` (
`id` varchar(191) NOT NULL DEFAULT '',
`user_id` varchar(191) NOT NULL,
`name` varchar(191) NOT NULL DEFAULT '',
`location` varchar(191) NOT NULL DEFAULT '',
`description` text NOT NULL,
`size` int(11) NOT NULL,
`created_at` datetime NOT NULL,
PRIMARY KEY (`id`),
KEY `user_id_idx` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE `sessions` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`session_name` VARCHAR(191) NOT NULL DEFAULT '',
`session_id` varchar(191) NOT NULL DEFAULT '',
`user_id` VARCHAR(191) NOT NULL DEFAULT '',
`session_expiry` VARCHAR(191) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `session_name` (`session_name`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE `users` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`user_id` VARCHAR(191) NOT NULL DEFAULT '',
`user_email` VARCHAR(128) NOT NULL DEFAULT '',
`user_password` VARCHAR(191) NOT NULL DEFAULT '',
`user_name` VARCHAR(191) NOT NULL DEFAULT '',
`user_level` int(11) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE `articles` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `page_user_id` varchar(191) NOT NULL,
  `page_guid` varchar(191) NOT NULL DEFAULT '',
  `page_title` varchar(191) DEFAULT NULL,
  `page_image_id` varchar(191) DEFAULT NULL,
  `page_content` mediumtext,
  `page_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON
  UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `page_guid` (`page_guid`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
