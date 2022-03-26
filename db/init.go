package db

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github-cve/config"
)

var DB *gorm.DB

func init() {
	//用户名:密码@tcp(IP:端口)/数据库?charset=utf8
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		config.Config.Username, config.Config.Password, config.Config.Host, config.Config.Port))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		panic(err)
	}
	_, err = db.Exec(fmt.Sprintf("Create Database If Not Exists %s Character Set UTF8", "cve"))
	if err != nil {
		panic(err)
	}
	//user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.Username, config.Config.Password, config.Config.Host, config.Config.Port, config.Config.Dbname)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // 表名前缀，`User`表为`t_users`
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	gormDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	err = gormDB.Ping()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	gormDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	gormDB.SetMaxOpenConns(100)
	log.Println("连接数据库成功")
}
