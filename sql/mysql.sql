-- 创建数据库 aphrodite
CREATE DATABASE `aphrodite` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;

-- stock 表
drop table if exists `aphrodite`.`stock`;
create table `aphrodite`.`stock` (
    `source` VARCHAR(8) NOT NULL COMMENT '来源',
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `name` VARCHAR(32) NOT NULL COMMENT '名称',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY(`code`)
);


-- quote_day 表
drop table if exists `aphrodite`.`quote_day`;
create table `aphrodite`.`quote_day` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `open` DECIMAL NOT NULL COMMENT '开盘价',
    `close` DECIMAL NOT NULL COMMENT '收盘价',
    `high` DECIMAL NOT NULL COMMENT '最高价',
    `low` DECIMAL NOT NULL COMMENT '最低价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL NOT NULL COMMENT '金额',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `year_day` INT NOT NULL COMMENT '天数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);

-- stock 表
drop table if exists `aphrodite`.`quote_week`;
create table `aphrodite`.`quote_week` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `open` DECIMAL NOT NULL COMMENT '开盘价',
    `close` DECIMAL NOT NULL COMMENT '收盘价',
    `high` DECIMAL NOT NULL COMMENT '最高价',
    `low` DECIMAL NOT NULL COMMENT '最低价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL NOT NULL COMMENT '金额',
    `date_begin` TIMESTAMP NOT NULL COMMENT '开始时期',
    `date_end` TIMESTAMP NOT NULL COMMENT '结束时期',
    `year_week` INT NOT NULL COMMENT '周数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);