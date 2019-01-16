package main

import (
	_ "edgea/routers"
	_ "edgea/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
