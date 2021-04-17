# db

## init
### 1. create database
```sql
sudo mysql -u root
set password for 'root'@'localhost' = PASSWORD('xxx');

create user 'authkey'@'%' identified by 'authkey';
create database authkey character set utf8mb4 collate utf8mb4_unicode_ci;
grant all privileges on authkey.* to 'authkey'@'%';
flush privileges;
```

### 2. init data
```sql
create table `user` (
    `id` bigint(20) not null auto_increment,
    `name` varchar(64) not null default '',
    `mobile` varchar(64) not null default '',
    `password` varchar(64) not null default '',
    `created_at` datetime not null default current_timestamp(),
    `properties` json default '{}',
    primary key (`id`),
    key (`mobile`)
) engine=innodb;
```