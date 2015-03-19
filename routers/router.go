package routers

import (
	"github.com/astaxie/beego"
	"zenapi/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/messageevent", &controllers.MessageEventController{})
	beego.Router("/messageevent/joinclient", &controllers.MessageEventController{}, "get:JoinClient")
	beego.Router("/zenapi/signal", &controllers.SignalController{}, "post:AddSignal")
	beego.Router("/s/s.gif", &controllers.SignalImageController{}, "get:GetSignalImage")
}
