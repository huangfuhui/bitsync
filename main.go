package main

import (
	"github.com/astaxie/beego"
	_ "bitsync/models"
	_ "bitsync/routers"
)

func main() {
	beego.Run()
}
