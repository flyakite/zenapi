package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	"strconv"
	. "zenapi/models"
	"zenapi/zen"
	//_ "github.com/lib/pq" //posgresql
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	zen.InitEnv()
	orm.RegisterModel(new(Access), new(Link), new(Receiver), new(Setting),
		new(Signal), new(User), new(UserAgent), new(UserTrack))
	zen.SyncDB()
}

func main() {

	//web
	port := beego.AppConfig.String("HttpPort")
	if port == "" {
		port, err := strconv.Atoi(os.Getenv("PORT")) //for heroku
		if err == nil {
			beego.HttpPort = port
		}
	}
	beego.Run()
}
