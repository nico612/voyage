-- Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file. The original repo for
-- this file is https://github.com/marmotedu/miniblog.


# 创建脚本，创建完成后最好使用`mysqldump`工具导出数据和表的创建SQL语句
DROP DATABASE IF EXISTS `miniblog`;
CREATE DATABASE `miniblog`;

USE miniblog;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`        bigint unsigned NOT NULL AUTO_INCREMENT,
    `username`  varchar(255)    NOT NULL,
    `password`  varchar(255)    NOT NULL,
    `nickname`  varchar(30)     NOT NULL,
    `email`     varchar(256)    NOT NULL,
    `phone`     varchar(16)     NOT NULL,
    `createdAt` timestamp       NOT NULL DEFAULT current_timestamp(),
    `updatedAt` timestamp       NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 2
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `post`;
CREATE TABLE `post`
(
    `id`        bigint unsigned NOT NULL AUTO_INCREMENT,
    `username`  varchar(255)    NOT NULL,
    `postID`    varchar(256)    NOT NULL,
    `title`     varchar(256)    NOT NULL,
    `content`   longtext        NOT NULL,
    `createdAt` timestamp       NOT NULL DEFAULT current_timestamp(),
    `updatedAt` timestamp       NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `postID` (`postID`),
    KEY `idx_username` (`username`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 2
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;