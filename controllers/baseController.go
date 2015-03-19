package controllers

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

const base64PixelGif = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

func (c *BaseController) ServePixelImage() {
	c.Ctx.Output.ContentType("image/gif")
	c.Ctx.Output.Header("cache-control", "private, max-age=0, no-cache")
	c.Ctx.Output.Body(base64PixelGif)
}
