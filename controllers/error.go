package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	c.TplName = "error/404.html"
}

func (c *ErrorController) Error500() {
	c.TplName = "error/500.html"
}
