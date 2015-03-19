package zen

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"os"
	"time"
	//_ "github.com/lib/pq" //posgresql
	_ "github.com/mattn/go-sqlite3"
)

var Cache cache.Cache
var RunMode string
var Config = beego.AppConfig

func InitEnv() {
	var err error
	beego.Debug("InitEnv")
	beego.SetLogger("file", `{"filename":"logs/debug.log"}`) //log to both console and file
	beego.SetLogFuncCall(true)                               //output file and line number
	beego.SetLevel(beego.LevelInformational)

	//RunMode
	RunMode = Config.String("runmode")

	//cache
	Cache, err := cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		beego.Critical("Cache Init Failed")
		os.Exit(1)
	}
	beego.Debug(Cache)

	//database
	maxIdle := 30
	maxConn := 30
	orm.RegisterDriver("sqlite3", orm.DR_Sqlite)
	orm.RegisterDataBase("default", "sqlite3", "data.db", maxIdle, maxConn)
	orm.DefaultTimeLoc = time.UTC

	if RunMode == "dev" {
		orm.Debug = true
	}

}

func SyncDB() {
	beego.Debug("syncDB")
	name := "default"
	force := false
	verbose := true
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		beego.Critical("Database Table Generation Failed")
		os.Exit(1)
	}
}
