package main

import (
	_ "api-generated/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

