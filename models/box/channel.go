package box

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrChannelNotFound     = errors.New("channel not found")
	ErrChannelAlreadyExist = errors.New("channel already exist")
)

type Channel struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 格子ID
	GridId int `orm:"column(grid_id)" json:"gridId"`
	// 传感器ID
	SensorId int `orm:"column(sensor_id)" json:"sensorId"`
	// 通道ID
	Channel int `orm:"column(channel)" json:"channel"`
	// 高度
	// TODO:自动触发测量高度
	Height int `orm:"column(height)" json:"height"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`

	// 传感器名称
	SensorName string `json:"sensorName"`
}

func (t *Channel) TableName() string {
	return "rel_grid_channel"
}

// 添加
func InsertChannel(obj *Channel) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelChannel(id int) error {
	o := orm.NewOrm()

	obj := &Channel{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateChannel(obj *Channel) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func ChannelById(id int) (*Channel, error) {
	o := orm.NewOrm()

	obj := &Channel{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrChannelNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据格子传感器ID查询
func ChannelByGridId(gridId, sensorId int) (*Channel, error) {
	o := orm.NewOrm()

	obj := &Channel{
		GridId:   gridId,
		SensorId: sensorId,
	}

	if err := o.Read(obj, "GridId", "SensorId"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrChannelNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func ChannelList(where map[string]interface{}, page, pageSize int) (int64, []*Channel, error) {
	o := orm.NewOrm()

	list := []*Channel{}

	sql := " 1 "
	if len(where) > 0 {
		startDate := where["startDate"]
		if startDate != "" {
			sql += " AND t1.created >= '" + fmt.Sprintf("%s", startDate) + "' "
		}

		endDate := where["endDate"]
		if endDate != "" {
			sql += " AND t1.created >= '" + fmt.Sprintf("%s", endDate) + "' "
		}

		gridId := where["gridId"]
		if gridId.(int) > 0 {
			sql += " AND t1.grid_id = " + fmt.Sprintf("%d", gridId) + " "
		}

		sensorId := where["sensorId"]
		if sensorId.(int) > 0 {
			sql += " AND t1.sensor_id = " + fmt.Sprintf("%d", sensorId) + " "
		}

		name := where["name"]
		if name != "" {
			sql += " AND (t1.name LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(channelListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(channelListSql+sql+" ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const channelListCountSql = `
SELECT
    COUNT(*)
FROM
    rel_grid_channel AS t1
WHERE
`

const channelListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.grid_id,
    t1.sensor_id,
    t1.channel,
    t1.height,
    t1.updated,
    t1.updated_by
FROM
    rel_grid_channel AS t1
WHERE
`

// 查询所有传感器
func SensorByGrid(gridId int) ([]*Channel, error) {
	o := orm.NewOrm()

	list := []*Channel{}

	// 查询所有
	if _, err := o.Raw(sensorByGridSql, gridId).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const sensorByGridSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.grid_id,
    t1.sensor_id,
    t1.channel,
    t1.height,
    t1.updated,
    t1.updated_by,
    t2.name AS sensor_name
FROM
    rel_grid_channel AS t1
INNER JOIN
    sensor AS t2
ON
    t1.sensor_id = t2.id
WHERE
    t1.grid_id = ?
`
