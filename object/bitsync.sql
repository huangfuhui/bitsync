-- 建立数据库
create database if not exists bitsync default CHARSET utf8 collate utf8_general_ci;

-- 模板
create table table_name (
  id int(10) unsigned not null auto_increment,
  uid int(10) unsigned unique not null comment '用户UID',
  created_at timestamp null default null,
  updated_at timestamp null default null,
  primary key (id)
) engine=InnoDB default charset=utf8 comment='';

create table account (
  id int(10) unsigned not null auto_increment,
  uid int(10) unsigned unique not null comment '用户UID',
  account varchar(20) not null comment '账号',
  password varchar(256) not null comment '密码',
  wechat_openid varchar(50) not null comment '微信openid',
  status tinyint(3) unsigned default '1' not null comment '账号状态(1:可用 2:停用)',
  register_time timestamp null comment '注册时间',
  login_time timestamp null comment '本次登录时间',
  login_ip varchar(15) null comment '本次登录IP',
  last_login_time timestamp null comment '上一次登录时间',
  last_login_ip varchar(15) null comment '上一次登录IP',
  created_at timestamp null default null,
  updated_at timestamp null default null,
  primary key (id),
  key idx_account_password (account, password)
) engine=InnoDB default charset=utf8 comment='账户';

create table member (
  id int(10) unsigned not null auto_increment,
  uid int(10) unsigned not null comment '用户UID',
  name varchar(20) not null comment '昵称',
  handset varchar(20) not null comment '手机号码',
  email varchar(50) not null comment '邮箱',
  sex tinyint(3) unsigned default '1' not null comment '性别(0:女 1:男)',
  avatar_url text null comment '头像URL',
  birthday date null comment '生日',
  created_at timestamp null default null,
  updated_at timestamp null default null,
  primary key (id)
) engine=InnoDB default charset=utf8 comment='会员';

create table sms_wallet (
  id int(10) unsigned not null auto_increment,
  uid int(10) unsigned unique not null comment '用户UID',
  balance int(10) unsigned not null default '0' comment '余额',
  prepare_consume int(10) unsigned not null default '0' comment '预消费',
  created_at timestamp null default null,
  updated_at timestamp null default null,
  primary key (id)
) engine=InnoDB default charset=utf8 comment='短信钱包';

create table sms_consume_record (
  id int(10) unsigned not null auto_increment,
  uid int(10) unsigned unique not null comment '用户UID',
  handset varchar(20) not null comment '手机号码',
  amount int(10) unsigned not null comment '消费条数',
  sms_content text not null comment '短信内容',
  consume_time timestamp not null comment '消费时间',
  status tinyint(3) unsigned not null default '1' comment '消费状态 [1.消费成功 2.消费失败]',
  created_at timestamp null default null,
  updated_at timestamp null default null,
  primary key (id)
) engine=InnoDB default charset=utf8 comment='短信消费记录';

create table sms_task (
  id int(10) unsigned not null auto_increment,
  uid int(10) unsigned unique not null comment '用户UID',
  task_id int(10) unsigned not null comment '任务ID',
  type tinyint(3) unsigned not null comment '任务类型 [1.阈值提醒]',
  status tinyint(3) unsigned not null default '0' comment '任务状态 [0.等待 1.成功 2.失败 3.取消]',
  finish_time timestamp null default null comment '任务完成时间',
  created_at timestamp null default null,
  updated_at timestamp null default null,
  primary key (id)
) engine=InnoDB default charset=utf8 comment='短信任务';

create table sms_failed_task (
  id int(10) unsigned not null auto_increment,
  uid int(10) unsigned unique not null comment '用户UID',
  sms_task_id int(10) unsigned not null comment '短信任务ID',
  failed_reason varchar(100) default '' comment '失败原因',
  created_at timestamp null default null,
  updated_at timestamp null default null,
  primary key (id)
) engine=InnoDB default charset=utf8 comment='短信失败任务';

create table sms_task_threshold_value (
  id int(10) unsigned not null auto_increment,
  uid int(10) unsigned unique not null comment '用户UID',
  symbol_pari varchar(20) not null comment '价格对',
  exchange tinyint(3) unsigned not null comment '交易所 [1.huobi 2.dragonex]',
  threshold_value varchar(20) not null comment '阈值',
  base_vale varchar(20) not null comment '基准值',
  deviation tinyint(1) not null comment '偏离方向 [1.大于 2.小于]',
  created_at timestamp null default null,
  updated_at timestamp null default null,
  primary key (id)
) engine=InnoDB default charset=utf8 comment='短信阈值提醒任务';
