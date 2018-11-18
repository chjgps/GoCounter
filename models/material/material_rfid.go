package material

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrMaterialRfidNotFound     = errors.New("material rfid not found")
	ErrMaterialRfidAlreadyExist = errors.New("material rfid already exist")
)

type MaterialRfidObj struct {
	*MaterialRfid
	Rfids []string `json:"rfids"`
}

type MaterialRfid struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 物料ID
	MaterialId int `orm:"column(material_id)" json:"materialId"` // OneToMany
	// Rfid
	Rfid string `orm:"column(rfid);unique" json:"rfid"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *MaterialRfid) TableName() string {
	return "rel_material_rfid"
}

// 添加
func InsertMaterialRfid(obj *MaterialRfid) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelMaterialRfid(id int) error {
	o := orm.NewOrm()

	obj := &MaterialRfid{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateMaterialRfid(obj *MaterialRfid) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func MaterialRfidById(id int) (*MaterialRfid, error) {
	o := orm.NewOrm()

	obj := &MaterialRfid{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrMaterialRfidNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

/*
// 根据ID查询
func MaterialRfidById(materialId int, rfid string) (*MaterialRfid, error) {
	o := orm.NewOrm()

	obj := &MaterialRfid{
		MaterialId: materialId,
		Rfid:       rfid,
	}

	if err := o.Read(obj, "MaterialId", "Rfid"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrMaterialRfidNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}
*/

// 根据Rfid查询
func MaterialByRfid(rfid string) (*Material, error) {
	o := orm.NewOrm()

	obj := &Material{}

	if err := o.Raw(materialByRfidSql, rfid).QueryRow(obj); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrMaterialRfidNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

const materialByRfidSql = `
SELECT
    *
FROM
    material
WHERE
    id = (
        SELECT
            material_id
        FROM
            rel_material_rfid
        WHERE
            rfid = ?
    )
`

// 查询所有
func MaterialRfidList(where map[string]interface{}, page, pageSize int) (int64, []*MaterialRfid, error) {
	o := orm.NewOrm()

	list := []*MaterialRfid{}

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
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(materialRfidListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(materialRfidListSql+sql+" GROUP BY t1.id ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const materialRfidListCountSql = `
SELECT
    COUNT(*)
FROM
    rel_material_rfid AS t1
WHERE
`

const materialRfidListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.material_id,
    t1.rfid,
    t1.status,
    t1.updated,
    t1.updated_by
FROM
    rel_material_rfid AS t1
WHERE
`
