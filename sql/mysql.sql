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
    `yesterday_closed` DECIMAL(10,2) NOT NULL COMMENT '昨日收盘价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL(18,2) NOT NULL COMMENT '金额',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `day_of_year` INT NOT NULL COMMENT '天数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_code on `aphrodite`.`quote_day`(`code`,`date`);

-- quote_week 表
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
create index idx_date_end_code on `aphrodite`.`quote_week`(`code`,`date_end`);

-- task_record 表
drop table if exists `aphrodite`.`task_record`;
create table `aphrodite`.`task_record` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `method` VARCHAR(32) NOT NULL COMMENT '任务名称',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `priority` INT NOT NULL COMMENT '优先级',
    `completed` BOOLEAN NOT NULL COMMENT '是否完成',
    `num_of_times` INT NOT NULL COMMENT '次数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_completed on `aphrodite`.`task_record`(`date`,`completed`);

-- ma_day 表
drop table if exists `aphrodite`.`ma_day`;
create table `aphrodite`.`ma_day` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `m5` DECIMAL(10,2) NOT NULL COMMENT '5 日移动平均线',
    `m10` DECIMAL(10,2) NOT NULL COMMENT '10 日移动平均线',
    `m20` DECIMAL(10,2) NOT NULL COMMENT '20 日移动平均线',
    `m30` DECIMAL(10,2) NOT NULL COMMENT '30 日移动平均线',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `day_of_year` INT NOT NULL COMMENT '天数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_code on `aphrodite`.`ma_day`(`code`,`date`);


-- ma_week 表
drop table if exists `aphrodite`.`ma_week`;
create table `aphrodite`.`ma_week` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `m5` DECIMAL(10,2) NOT NULL COMMENT '5 周移动平均线',
    `m10` DECIMAL(10,2) NOT NULL COMMENT '10 周移动平均线',
    `m20` DECIMAL(10,2) NOT NULL COMMENT '20 周移动平均线',
    `m30` DECIMAL(10,2) NOT NULL COMMENT '30 周移动平均线',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `week_of_year` INT NOT NULL COMMENT '周数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_code on `aphrodite`.`ma_week`(`code`,`date`);


-- boll_day 表
drop table if exists `aphrodite`.`boll_day`;
create table `aphrodite`.`boll_day` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `up` DECIMAL(10,2) NOT NULL COMMENT '上轨线',
    `mb` DECIMAL(10,2) NOT NULL COMMENT '中轨线',
    `dn` DECIMAL(10,2) NOT NULL COMMENT '下轨线',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `day_of_year` INT NOT NULL COMMENT '天数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_code on `aphrodite`.`boll_day`(`code`,`date`);


-- boll_week 表
drop table if exists `aphrodite`.`boll_week`;
create table `aphrodite`.`boll_week` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `up` DECIMAL(10,2) NOT NULL COMMENT '上轨线',
    `mb` DECIMAL(10,2) NOT NULL COMMENT '中轨线',
    `dn` DECIMAL(10,2) NOT NULL COMMENT '下轨线',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `week_of_year` INT NOT NULL COMMENT '周数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_code on `aphrodite`.`boll_week`(`code`,`date`);


-- kdj_day 表
drop table if exists `aphrodite`.`kdj_day`;
create table `aphrodite`.`kdj_day` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `k` DECIMAL(10,2) NOT NULL COMMENT 'k 值',
    `d` DECIMAL(10,2) NOT NULL COMMENT 'd 值',
    `j` DECIMAL(10,2) NOT NULL COMMENT 'j 值',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `day_of_year` INT NOT NULL COMMENT '天数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_code on `aphrodite`.`kdj_day`(`code`,`date`);

-- kdj_week 表
drop table if exists `aphrodite`.`kdj_week`;
create table `aphrodite`.`kdj_week` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `k` DECIMAL(10,2) NOT NULL COMMENT 'k 值',
    `d` DECIMAL(10,2) NOT NULL COMMENT 'd 值',
    `j` DECIMAL(10,2) NOT NULL COMMENT 'j 值',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `week_of_year` INT NOT NULL COMMENT '周数',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_date_code on `aphrodite`.`kdj_week`(`code`,`date`);
