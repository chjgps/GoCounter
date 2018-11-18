package material

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryAlreadyExist = errors.New("category already exist")
)

type Category struct {
	// 类别ID
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 父类ID
	ParentId string `orm:"column(parent_id)" json:"parentId"`
	// 名称
	Name string `orm:"column(name)" json:"name"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *Category) TableName() string {
	return "rel_material_category"
}

// 添加
func InsertCategory(obj *Category) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateCategory(obj *Category) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelCategory(id int) error {
	o := orm.NewOrm()

	obj := &Category{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func CategoryById(id int) (*Category, error) {
	o := orm.NewOrm()

	obj := &Category{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrCategoryNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据名称查询
func CategoryByName(name string) (*Category, error) {
	o := orm.NewOrm()

	obj := &Category{
		Name: name,
	}

	if err := o.Read(obj, "Name"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrCategoryNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询分类
func CategoryList() ([]*Category, error) {
	o := orm.NewOrm()

	list := []*Category{}
	if _, err := o.Raw(categoryListSql).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const categoryListSql = `
WITH RECURSIVE
MATERIAL_CTE(id, parent_id, name, level) AS (
	SELECT
		t1.id,
		t1.parent_id,
		t1.name,
		0 AS level
	FROM
		rel_material_category AS t1
	WHERE
		t1.parent_id = 0
	UNION ALL
		SELECT
			t2.id,
			t2.parent_id,
			t2.name,
			t3.level + 1
		FROM
			rel_material_category AS t2
		INNER JOIN MATERIAL_CTE AS t3 ON t2.parent_id = t3.id
)
SELECT
	*
FROM
	MATERIAL_CTE;
`
