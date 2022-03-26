# github-cve

github cve monitor

## 安装mariadb

<span style="color: red">注意：mariadb版本>5.5</span>

### 创建yum源文件

```shell
vim /etc/yum.repos.d/mariaDB.repo
```

### 输入内容

```shell
# MariaDB 10.7 CentOS repository list - created 2022-03-26 10:12 UTC
# https://mariadb.org/download/
[mariadb]
name = MariaDB
baseurl = https://mirrors.aliyun.com/mariadb/yum/10.7/centos7-amd64
gpgkey=https://mirrors.aliyun.com/mariadb/yum/RPM-GPG-KEY-MariaDB
gpgcheck=1
```

### 安装mariadb-server

```shell
yum install MariaDB-server 
```

## 配置数据库

config.ini

```shell
[mysql]
host = localhost # 数据库地址
port = 3306      # 数据库端口
dbname = cve     # 数据库名称
username = root  # 用户名
password =       # 密码
```

## 启动

```shell
go run main.go
```