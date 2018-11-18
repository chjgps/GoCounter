package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/beego/ms304w-client/basis/conf"
	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/httpx/rest"
)

type SerialController struct {
	BaseController
}

const (
	LOCK_UNWEIGHT    = 0           // 不称重
	LOCK_WEIGHT      = 1           // 称重
	LOCK_SUM         = 2           // 计数
	LOCK_CHECK       = 3           // 盘点
	WEIGHT_UNDEFINED = "undefined" // UUID不存在，只称重，不处理数据
	// 灯
	LIGHT_OPEN  int = 1
	LIGHT_CLOSE int = 0
)

var (
	// API路径
	API_URL = conf.String("api_server")

	SERIAL_OPEN_URL = API_URL + "/v1/serial/open"

	SERIAL_WEIGHT_URL = API_URL + "/v1/serial/weight"

	SERIAL_WEIGHT_ZERO = API_URL + "/v1/serial/weight/zero"

	SERIAL_WEIGHT_MEASURE = API_URL + "/v1/serial/weight/measure"

	SERIAL_WEIGHT_CHECK = API_URL + "/v1/serial/weight/check"

	SERIAL_LIGHT = API_URL + "/v1/serial/light"

	SERIAL_ALL_STATUS = API_URL + "/v1/serial/box/%d/status"
)

type ApiData struct {
	// 业务唯一ID
	UUID string `json:"uuid"`
	// 柜子通信ID
	BoxId int `json:"boxId"`
	// 格子通道ID
	GridId int `json:"gridId"`
	// 操作类型
	// 0不称重1称重2计数3盘点
	Operation int `json:"operation"`
	// 重量,盘点
	Weight int `json:"weight"`
	// 门状态列表
	DoorStatusList string `json:"doorStatusList"`
	// 门状态
	DoorStatus int `json:"doorStatus"`
	// 灯状态
	LightStatus int `json:"lightStatus"` // 卡号
	// 卡号
	Card string `json:"card"`
	// 条码
	Code string `json:"code"`
	// 指纹
	Finger int `json:"finger"`
}

// 开门
func (c *SerialController) Open() {
	log.Info("Open: %s", string(c.Ctx.Input.RequestBody))

	obj := &ApiData{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	// boxId
	boxId := obj.BoxId
	if boxId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("boxId is empty"))
		return
	}

	// gridId
	gridId := obj.GridId
	if gridId < 0 {
		c.WriteHttpResponse(400, nil, errors.New("gridId is empty"))
		return
	}

	// operation
	operation := obj.Operation

	if operation != LOCK_WEIGHT && operation != LOCK_UNWEIGHT && operation != LOCK_SUM {
		c.WriteHttpResponse(400, nil, errors.New("operation is illegal").As(operation))
		return
	}

	where := make(map[string]interface{})
	where["boxId"] = boxId
	where["gridId"] = gridId
	where["operation"] = operation

	log.Info("Post %s", SERIAL_OPEN_URL)

	bytes, err := rest.Post(SERIAL_OPEN_URL).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("Open Res %s", string(bytes))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 称重
func (c *SerialController) Weight() {
	log.Info("Weight: %s", string(c.Ctx.Input.RequestBody))

	obj := &ApiData{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	// boxId
	boxId := obj.BoxId
	if boxId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("boxId is empty"))
		return
	}

	// gridId
	gridId := obj.GridId
	if gridId < 0 {
		c.WriteHttpResponse(400, nil, errors.New("gridId is empty"))
		return
	}

	where := make(map[string]interface{})
	where["boxId"] = boxId
	where["gridId"] = gridId

	bytes, err := rest.Post(SERIAL_WEIGHT_URL).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("Weight Res %s", string(bytes))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 清零
func (c *SerialController) Zero() {
	log.Info("Weight Zero: %s", string(c.Ctx.Input.RequestBody))

	obj := &ApiData{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	// boxId
	boxId := obj.BoxId
	if boxId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("boxId is empty"))
		return
	}

	// gridId
	gridId := obj.GridId
	if gridId < 0 {
		c.WriteHttpResponse(400, nil, errors.New("gridId is empty"))
		return
	}

	where := make(map[string]interface{})
	where["boxId"] = boxId
	where["gridId"] = gridId

	bytes, err := rest.Post(SERIAL_WEIGHT_ZERO).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("Zero Res %s", string(bytes))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 砝码矫正
func (c *SerialController) Measure() {
	log.Info("Weight Measure: %s", string(c.Ctx.Input.RequestBody))

	obj := &ApiData{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	// boxId
	boxId := obj.BoxId
	if boxId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("boxId is empty"))
		return
	}

	// gridId
	gridId := obj.GridId
	if gridId < 0 {
		c.WriteHttpResponse(400, nil, errors.New("gridId is empty"))
		return
	}

	// weight
	weight := obj.Weight
	if weight <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("weight is empty"))
		return
	}

	where := make(map[string]interface{})
	where["boxId"] = boxId
	where["gridId"] = gridId
	where["weight"] = weight

	bytes, err := rest.Post(SERIAL_WEIGHT_MEASURE).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("Measure Res %s", string(bytes))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 盘点
func (c *SerialController) Check() {
	log.Info("Check: %s", string(c.Ctx.Input.RequestBody))

	obj := &ApiData{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	// uuid
	uuid := obj.UUID
	if len(uuid) <= 0 {
		uuid = WEIGHT_UNDEFINED
	}

	// boxId
	boxId := obj.BoxId
	if boxId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("boxId is empty"))
		return
	}

	// gridId
	gridId := obj.GridId
	if gridId < 0 {
		c.WriteHttpResponse(400, nil, errors.New("gridId is empty"))
		return
	}

	where := make(map[string]interface{})
	where["boxId"] = boxId
	where["gridId"] = gridId

	bytes, err := rest.Post(SERIAL_WEIGHT_CHECK).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("Check Res %s", string(bytes))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 获取所有门状态
func (c *SerialController) BoxStatus() {
	log.Info("All Door Status: %s", string(c.Ctx.Input.RequestBody))

	boxIdStr := c.Ctx.Input.Param(":id")
	log.Debug(boxIdStr)
	if len(boxIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("box id is empty"))
		return
	}

	// boxId
	boxId, err := strconv.Atoi(boxIdStr)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if boxId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("boxId is empty"))
		return
	}

	bytes, err := rest.Get(fmt.Sprintf(SERIAL_ALL_STATUS, boxId)).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("All Status Res %s", string(bytes))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 开关灯
func (c *SerialController) Light() {
	log.Info("Light: %s", string(c.Ctx.Input.RequestBody))

	obj := &ApiData{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	// boxId
	boxId := obj.BoxId
	if boxId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("boxId is empty"))
		return
	}

	// gridId
	gridId := obj.GridId
	if gridId < 0 {
		c.WriteHttpResponse(400, nil, errors.New("gridId is empty"))
		return
	}

	// operation
	operation := obj.Operation
	if operation != LIGHT_OPEN && operation != LIGHT_CLOSE {
		c.WriteHttpResponse(400, nil, errors.New("operation is illegal").As(operation))
		return
	}

	where := make(map[string]interface{})
	where["boxId"] = boxId
	where["gridId"] = gridId
	where["operation"] = operation

	bytes, err := rest.Post(SERIAL_LIGHT).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("Light Res %s", string(bytes))

	c.WriteHttpResponse(200, nil, nil)
	return
}
