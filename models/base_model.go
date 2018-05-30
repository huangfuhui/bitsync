package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BaseModel struct {
}

func init() {
	mysqlUser := beego.AppConfig.String("mysql_user")
	mysqlPassword := beego.AppConfig.String("mysql_password")
	mysqlHost := beego.AppConfig.String("mysql_host")
	mysqlPort := beego.AppConfig.String("mysql_port")
	mysqlDb := beego.AppConfig.String("mysql_db")

	dataSource := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDb + "?charset=utf8"
	// 注册数据库
	orm.RegisterDataBase("defalut", "mysql", dataSource)
}
