package test

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "zenapi/models"
	"zenapi/zen"
)

func testInit() {
	beego.Debug("testInit")
	zen.InitEnv()
	orm.RegisterModel(new(Access), new(Link), new(Receiver), new(Setting),
		new(Signal), new(User), new(UserAgent), new(UserTrack))
	zen.SyncDB()
}

func init() {
	testInit()
}
