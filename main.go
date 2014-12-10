package main

import (
	"github.com/astaxie/beego"
	"os"
	"strconv"
	_ "zenapi/routers"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		beego.HttpPort = port
	}

	beego.Run()
}
