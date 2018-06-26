package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"bitsync/object/member"
)

type BaseModel struct {
}

func init() {
	runMode := beego.AppConfig.String("runmode")

	if runMode != "prod" {
		orm.Debug = true
	}

	mysqlUser := beego.AppConfig.String("mysql_user")
	mysqlPassword := beego.AppConfig.String("mysql_password")
	mysqlHost := beego.AppConfig.String("mysql_host")
	mysqlPort := beego.AppConfig.String("mysql_port")
	mysqlDb := beego.AppConfig.String("mysql_db")

	dataSource := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDb + "?charset=utf8"
	// 注册数据库
	err := orm.RegisterDataBase("default", "mysql", dataSource)
	if err != nil {
		beego.Error(err)
	}

	// 注册模型
	orm.RegisterModel(
		new(member.Account),
		new(member.Member),
	)
}
