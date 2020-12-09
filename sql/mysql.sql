-- 创建数据库 aphrodite
CREATE DATABASE `aphrodite` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;

-- stock 表
drop table if exists `aphrodite`.`stock`;
create table `aphrodite`.`stock` (
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `name` VARCHAR(32) NOT NULL COMMENT '名称',
    `source` VARCHAR(8) NOT NULL COMMENT '来源',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY(`code`)
);


-- quote_day 表
drop table if exists `aphrodite`.`quote_day`;
create table `aphrodite`.`quote_day` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `open` DECIMAL(10,2) NOT NULL COMMENT '开盘价',
    `close` DECIMAL(10,2) NOT NULL COMMENT '收盘价',
    `high` DECIMAL(10,2) NOT NULL COMMENT '最高价',
    `low` DECIMAL(10,2) NOT NULL COMMENT '最低价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL(18,2) NOT NULL COMMENT '金额',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `day_of_year` INT NOT NULL COMMENT '天数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_code on `aphrodite`.`quote_day`(`date`,`code`);

-- stock 表
drop table if exists `aphrodite`.`quote_week`;
create table `aphrodite`.`quote_week` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `open` DECIMAL(10,2) NOT NULL COMMENT '开盘价',
    `close` DECIMAL(10,2) NOT NULL COMMENT '收盘价',
    `high` DECIMAL(10,2) NOT NULL COMMENT '最高价',
    `low` DECIMAL(10,2) NOT NULL COMMENT '最低价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL(18,2) NOT NULL COMMENT '金额',
    `date_begin` TIMESTAMP NOT NULL COMMENT '开始时期',
    `date_end` TIMESTAMP NOT NULL COMMENT '结束时期',
    `week_of_year` INT NOT NULL COMMENT '周数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_end_code on `aphrodite`.`quote_week`(`date_end`,`code`);

-- stock 表
drop table if exists `aphrodite`.`task_record`;
create table `aphrodite`.`task_record` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(32) NOT NULL COMMENT '任务名称',
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `completed` BOOLEAN NOT NULL COMMENT '是否完成',
    `msg` TEXT NOT NULL COMMENT '描述',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_completed on `aphrodite`.`task_record`(`date`,`completed`);