package box

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrGridNotFound     = errors.New("grid not found")
	ErrGridAlreadyExist = errors.New("grid already exist")
)

type Grid struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 名称
	Name string `orm:"column(name)" json:"name"`
	// 柜子ID
	BoxId int `orm:"column(box_id)" json:"boxId"`
	// 通道ID
	Channel int `orm:"column(channel)" json:"channel"`
	// 最大数量
	Qty int `orm:"column(qty)" json:"qty"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
	// 条码
	Code string `orm:"column(code)" json:"code"`
	// 物料
	MaterialId int `orm:"column(material_id)" json:"materialId"`
	// 类型(私有0,公有1)
	Type int `orm:"column(type)" json:"type"`
	// 安全库存
	SafeQty int `orm:"column(safe_qty)" json:"safeQty"`

	// other
	// box
	Addr int `json:"addr"`
	// 传感器
	SensorNameList string `json:"sensorNameList"`
	AccountGridId  int    `json:"accountGridId"`
	AccountId      int    `json:"accountId"`
	AccountName    string `json:"accountName"`
	// 分类ID
	CategoryId   int    `json:"categoryId"`
	MaterialCode string `json:"materialCode"`
	MaterialName string `json:"materialName"`
	TotalQty     int    `json:"totalQty"`
}

func (t *Grid) TableName() string {
	return "rel_box_grid"
}

// 添加
func InsertGrid(obj *Grid) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelGrid(id int) error {
	o := orm.NewOrm()

	obj := &Grid{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateGrid(obj *Grid) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func GridById(id int) (*Grid, error) {
	o := orm.NewOrm()

	obj := &Grid{}

	if err := o.Raw(gridByIdSql, id).QueryRow(obj); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrGridNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

const gridByIdSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.name,
    t1.box_id,
    t1.channel,
    t1.qty,
    t1.status,
    t1.updated,
    t1.updated_by,
    t1.code,
    t1.material_id,
    t1.type,
    t1.safe_qty,
    t2.addr AS addr,
    t3.category_id AS category_id,
    t3.material_code,
    SUM(t4.qty) AS total_qty
FROM
    rel_box_grid AS t1
LEFT JOIN
    box AS t2
ON
    t1.box_id = t2.id
LEFT JOIN
    material AS t3
ON
    t1.material_id = t3.id
LEFT JOIN
    stock AS t4
ON
    t1.material_id = t4.material_id
WHERE
    t1.id = ?
`

// 根据名称查询
func GridByName(name string) (*Grid, error) {
	o := orm.NewOrm()

	obj := &Grid{
		Name: name,
	}

	if err := o.Read(obj, "Name"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrGridNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据名称查询
func GridByMaterialId(materialId int) ([]*Grid, error) {
	o := orm.NewOrm()

	list := []*Grid{}

	if _, err := o.Raw(gridByMaterialIdSql, materialId).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const gridByMaterialIdSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.name,
    t1.box_id,
    t1.channel,
    t1.qty,
    t1.status,
    t1.updated,
    t1.updated_by,
    t1.code,
    t1.material_id,
    t1.type,
    t1.safe_qty,
    t2.addr AS addr,
    t3.category_id AS category_id,
    t3.material_code,
    SUM(t4.qty) AS total_qty
FROM
    rel_box_grid AS t1
LEFT JOIN
    box AS t2
ON
    t1.box_id = t2.id
LEFT JOIN
    material AS t3
ON
    t1.material_id = t3.id
LEFT JOIN
    stock AS t4
ON
    t1.material_id = t4.material_id
WHERE
    t1.material_id = ?
ORDER BY SUM(t4.qty) DESC
`

// 根据code查询
func GridByCode(code string) (*Grid, error) {
	o := orm.NewOrm()

	obj := &Grid{}

	if err := o.Raw(gridByCodeSql, code).QueryRow(obj); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrGridNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

const gridByCodeSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.name,
    t1.box_id,
    t1.channel,
    t1.qty,
    t1.status,
    t1.updated,
    t1.updated_by,
    t1.code,
    t1.material_id,
    t1.type,
    t1.safe_qty,
    t2.addr AS addr,
    t3.category_id AS category_id,
    t3.material_code,
    t3.name AS material_name,
    SUM(t4.qty) AS total_qty
FROM
    rel_box_grid AS t1
LEFT JOIN
    box AS t2
ON
    t1.box_id = t2.id
LEFT JOIN
    material AS t3
ON
    t1.material_id = t3.id
LEFT JOIN
    stock AS t4
ON
    t1.material_id = t4.material_id
WHERE
    t1.code = ?
`

// 查询所有
func GridList(where map[string]interface{}, page, pageSize int) (int64, []*Grid, error) {
	o := orm.NewOrm()

	list := []*Grid{}

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

		boxId := where["boxId"]
		if boxId.(int) > 0 {
			sql += " AND t1.box_id = " + fmt.Sprintf("%d", boxId) + " "
		}

		sensorId := where["sensorId"]
		if sensorId.(int) > 0 {
			sql += " AND t6.sensor_id = " + fmt.Sprintf("%d", sensorId) + " "
		}

		name := where["name"]
		if name != "" {
			sql += " AND (t1.name LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(gridListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(gridListSql+sql+" GROUP BY t1.id ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const gridListCountSql = `
SELECT
    COUNT(*)
FROM
    rel_box_grid AS t1
LEFT JOIN
    box AS t2
ON
    t1.box_id = t2.id
LEFT JOIN
    stock AS t3
ON
    t1.id = t3.grid_id
LEFT JOIN
    material AS t4
ON
    t1.material_id = t4.id
LEFT JOIN
    rel_account_grid AS t5
ON
    t1.id = t5.grid_id
LEFT JOIN
    account AS t6
ON
    t5.account_id = t6.id
LEFT JOIN
    rel_grid_channel AS t7
ON
    t1.id = t7.grid_id
WHERE
`

const gridListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.name,
    t1.box_id,
    t1.channel,
    t1.qty,
    t1.status,
    t1.updated,
    t1.updated_by,
    t1.code,
    t1.material_id,
    t1.type,
    t1.safe_qty,
    t2.addr AS addr,
    SUM(t3.qty) AS total_qty,
    t4.category_id AS category_id,
    t4.material_code,
    t5.id AS account_grid_id,
    t6.id AS account_id,
    t6.username AS account_name
FROM
    rel_box_grid AS t1
LEFT JOIN
    box AS t2
ON
    t1.box_id = t2.id
LEFT JOIN
    stock AS t3
ON
    t1.id = t3.grid_id
LEFT JOIN
    material AS t4
ON
    t1.material_id = t4.id
LEFT JOIN
    rel_account_grid AS t5
ON
    t1.id = t5.grid_id
LEFT JOIN
    account AS t6
ON
    t5.account_id = t6.id
LEFT JOIN
    rel_grid_channel AS t7
ON
    t1.id = t7.grid_id
WHERE
`
