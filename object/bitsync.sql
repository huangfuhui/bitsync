-- 建立数据库
create database if not exists bitsync default CHARSET utf8 collate utf8_general_ci;

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
