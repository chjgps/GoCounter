package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/beego/ms304w-client/basis/errors"
)

type BaseController struct {
	beego.Controller
}

/*
const httpResponseData = `
{
    "api": "1.0",
    "code": %d,
    "err": "%s",
    "data": "%s"
}
`
*/

type HttpResponse struct {
	Api  string      `json:"api"`
	Code int         `json:"code"`
	Err  string      `json:"err"`
	Data interface{} `json:"data"`
}

func NewHttpResponse(code int, data interface{}, err error) *HttpResponse {
	// data
	if data == nil {
		data = &struct{}{}
	}

	// error
	var errMsg string
	if err != nil {
		log.Error("%v", errors.As(err))
		errMsg = errors.ParseErr(err).Key()
	}

	return &HttpResponse{
		Api:  "1.0",
		Code: code,
		Err:  errMsg,
		Data: data,
	}
}

func (o *HttpResponse) String() string {
	bytes, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintf("%#v", err)
	}

	return string(bytes)
}

func (c *BaseController) WriteHttpResponse(code int, data interface{}, err error) {
	// data
	if data == nil {
		data = &struct{}{}
	}

	// error
	var errMsg string
	if err != nil {
		log.Error("%v", errors.As(err))
		errMsg = errors.ParseErr(err).Key()
	}

	// success
	c.Data["json"] = &HttpResponse{
		Api:  "1.0",
		Code: code,
		Err:  errMsg,
		Data: data,
	}

	c.ServeJSON()
	// c.StopRun()
}
