CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id字段',
  `name` varchar(100) NOT NULL COMMENT '名称',
  `age` int(11) NOT NULL DEFAULT 0 COMMENT '年龄',
  `ctime` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_name_IDX` (`name`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4