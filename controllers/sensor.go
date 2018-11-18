package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/sensor"
)

type SensorController struct {
	BaseController
}

// 添加
func (c *SensorController) AddSensor() {
	obj := &sensor.Sensor{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	name := obj.Name

	// 查询用户名是否重复
	acc, err := sensor.SensorByName(name)
	if err != nil {
		if !sensor.ErrSensorNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(sensor.ErrSensorAlreadyExist))
		return
	}

	obj.Status = 1
	obj.Created = timex.String()
	if err := sensor.InsertSensor(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改
func (c *SensorController) EditSensor() {
	obj := &sensor.Sensor{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	sensorId := obj.Id

	// 查询用户是否存在
	acc, err := sensor.SensorById(sensorId)
	if err != nil {
		if !sensor.ErrSensorNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(sensor.ErrSensorNotFound))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := sensor.UpdateSensor(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *SensorController) DelSensor() {
	sensorIdStr := c.Ctx.Input.Param(":id")
	log.Debug(sensorIdStr)
	if len(sensorIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("sensor id is empty"))
		return
	}

	sensorId, err := strconv.Atoi(sensorIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if err := sensor.DelSensor(sensorId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询
func (c *SensorController) SensorById() {
	sensorIdStr := c.Ctx.Input.Param(":id")
	log.Debug(sensorIdStr)
	if len(sensorIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("sensor id is empty"))
		return
	}

	sensorId, err := strconv.Atoi(sensorIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	acc, err := sensor.SensorById(sensorId)
	if err != nil {
		if !sensor.ErrSensorNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 查询所有
func (c *SensorController) SensorList() {
	page, err := c.GetInt("page")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	pageSize, err := c.GetInt("pageSize")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	list, err := sensor.SensorList(page, pageSize)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
