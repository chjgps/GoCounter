package material

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrSensorNotFound     = errors.New("material sensor not found")
	ErrSensorAlreadyExist = errors.New("material sensor already exist")
)

type Sensor struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 物料ID
	MaterialId int `orm:"column(material_id)" json:"materialId"` // OneToMany
	// 传感器
	SensorId int `orm:"column(sensor_id)" json:"sensorId"`
	// 参数
	/*
		        称重
			   {
			       "weight": 1,
			       "comeUp": 0.1,
			       "lower": 0.1
			   }
		       测距
			   {
			       "height": 1,
			       "comeUp": 0.1,
			       "lower": 0.1
			   }
	*/
	Params string `orm:"column(params)" json:"params"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`

	// other
	SensorName string `json:"sensorName"`
}

func (t *Sensor) TableName() string {
	return "rel_material_sensor"
}

type MaterialParams struct {
	Weight int `json:"weight"`
	Height int `json:"height"`
	ComeUp int `json:"comeUp"`
	Lower  int `json:"lower"`
}

func (t *Sensor) ParamsObj() (*MaterialParams, error) {
	obj := &MaterialParams{}

	if err := json.Unmarshal([]byte(t.Params), obj); err != nil {
		return nil, errors.As(err)
	}

	return obj, nil
}

// 添加
func InsertSensor(obj *Sensor) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelSensor(id int) error {
	o := orm.NewOrm()

	obj := &Sensor{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateSensor(obj *Sensor) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func SensorById(id int) (*Sensor, error) {
	o := orm.NewOrm()

	obj := &Sensor{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrSensorNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据物料和传感器ID查询
func SensorByMaterialId(materialId, sensorId int) (*Sensor, error) {
	o := orm.NewOrm()

	obj := &Sensor{
		MaterialId: materialId,
		SensorId:   sensorId,
	}

	if err := o.Read(obj, "MaterialId", "SensorId"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrSensorNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func SensorList(where map[string]interface{}, page, pageSize int) (int64, []*Sensor, error) {
	o := orm.NewOrm()

	list := []*Sensor{}

	sql := " 1 "
	if len(where) > 0 {
		startDate := where["startDate"]
		if startDate != "" {
			sql += " AND t1.created >= '" + fmt.Sprintf("%s", startDate) + "' "
		}

		endDate := where["endDate"]
		if endDate != "" {
			sql += " AND t1.created <= '" + fmt.Sprintf("%s", endDate) + "' "
		}

		materialId := where["materialId"]
		if materialId.(int) > 0 {
			sql += " AND t1.material_id = " + fmt.Sprintf("%d", materialId) + " "
		}

		name := where["name"]
		if name != "" {
			sql += " AND (t2.name LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(sensorListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(sensorListSql+sql+" ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const sensorListCountSql = `
SELECT
    COUNT(*)
FROM
    rel_material_sensor AS t1
WHERE
`

const sensorListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.material_id,
    t1.sensor_id,
    t1.params,
    t1.status,
    t1.updated,
    t1.updated_by,
    t2.name AS sensor_name
FROM
    rel_material_sensor AS t1
INNER JOIN
    sensor AS t2
ON
    t1.sensor_id = t2.id
WHERE
`
