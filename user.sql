CREATE TABLE IF NOT EXISTS  `user` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
`uid` varchar(45) NOT NULL COMMENT '玩家唯一ID',
`nickname` varchar(45) NOT NULL COMMENT '玩家昵称',
`icon` varchar(45) NOT NULL COMMENT '玩家头像',
`Forever` varchar(20000) NOT NULL COMMENT '永久存储',
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';