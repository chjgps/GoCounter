package material

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrMaterialNotFound     = errors.New("material not found")
	ErrMaterialAlreadyExist = errors.New("material already exist")
)

type Material struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 名称
	Name string `orm:"column(name)" json:"name"`
	// 分类ID
	CategoryId int `orm:"column(category_id)" json:"categoryId"` // OneToMany
	// 编码
	MaterialCode string `orm:"column(material_code)" json:"materialCode"`
	// 规格
	MaterialSpec string `orm:"column(material_spec)" json:"materialSpec"`
	// 供应商
	SupplierId int `orm:"column(supplier_id)" json:"supplierId"`
	// TODO:供应商库存
	SupplierQty int `orm:"column(supplier_qty)" json:"supplierQty"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
	// 图片地址
	Img string `orm:"column(img)" json:"img"`

	// other
	SupplierName string `json:"supplierName"`
	CategoryName string `json:"categoryName"`
	Qty          int    `json:"qty"`
}

func (t *Material) TableName() string {
	return "material"
}

// 添加
func InsertMaterial(obj *Material) (int, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return -1, errors.As(err)
	}

	return obj.Id, nil
}

// 删除
func DelMaterial(id int) error {
	o := orm.NewOrm()

	obj := &Material{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateMaterial(obj *Material) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func MaterialById(id int) (*Material, error) {
	o := orm.NewOrm()

	obj := &Material{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrMaterialNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据名称查询
func MaterialByName(name string) (*Material, error) {
	o := orm.NewOrm()

	obj := &Material{
		Name: name,
	}

	if err := o.Read(obj, "Name"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrMaterialNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func MaterialList(where map[string]interface{}, page, pageSize int) (int64, []*Material, error) {
	o := orm.NewOrm()

	list := []*Material{}

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

		categoryId := where["categoryId"]
		if categoryId.(int) > 0 {
			sql += " AND t1.category_id = " + fmt.Sprintf("%d", categoryId) + " "
		}

		name := where["name"]
		if name != "" {
			sql += " AND (t1.name LIKE '%" + fmt.Sprintf("%s", name) + "%'"
			sql += " OR t1.material_code LIKE '%" + fmt.Sprintf("%s", name) + "%'"
			sql += " OR t1.material_spec LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(materialListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(materialListSql+sql+" GROUP BY t1.id ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const materialListCountSql = `
SELECT
    COUNT(*)
FROM
    material AS t1
WHERE
`

const materialListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.name,
    t1.category_id,
    t1.material_code,
    t1.material_spec,
    t1.supplier_id,
    t1.supplier_qty,
    t1.status,
    t1.updated,
    t1.updated_by,
    t1.img,
    t2.short_name AS supplier_name,
    t3.name AS category_name,
    SUM(t4.qty) AS qty
FROM
    material AS t1
LEFT JOIN
    supplier AS t2
ON
    t1.supplier_id = t2.id
LEFT JOIN
    rel_material_category AS t3
ON
    t1.category_id = t3.id
LEFT JOIN
    stock AS t4
ON
    t1.id = t4.material_id
WHERE
`

// 查询所有
func SupplierMaterialList(where map[string]interface{}, page, pageSize int) (int64, []*Material, error) {
	o := orm.NewOrm()

	list := []*Material{}

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

		categoryId := where["categoryId"]
		if categoryId.(int) > 0 {
			sql += " AND t1.category_id = " + fmt.Sprintf("%d", categoryId) + " "
		}

		name := where["name"]
		if name != "" {
			sql += " AND (t1.name LIKE '%" + fmt.Sprintf("%s", name) + "%'"
			sql += " OR t1.material_code LIKE '%" + fmt.Sprintf("%s", name) + "%'"
			sql += " OR t1.material_spec LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(supplierMaterialListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(supplierMaterialListSql+sql+" GROUP BY t1.id ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const supplierMaterialListCountSql = `
SELECT
    COUNT(*)
FROM
    material AS t1
WHERE
`

const supplierMaterialListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.name,
    t1.category_id,
    t1.material_code,
    t1.material_spec,
    t1.supplier_id,
    t1.supplier_qty,
    t1.status,
    t1.updated,
    t1.updated_by,
    t1.img,
    t2.short_name AS supplier_name
FROM
    material AS t1
LEFT JOIN
    supplier AS t2
ON
    t1.supplier_id = t2.id
WHERE
`
