package controllers

import (
	l "github.com/beego/ms304w-client/basis/syslog"
	_ "github.com/beego/ms304w-client/models"
)

var log = l.New()

type MainController struct {
	BaseController
}

func (c *MainController) Index() {
	c.WriteHttpResponse(200, "ok", nil)
	return
}
