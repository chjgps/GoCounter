package material

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrSupplierNotFound     = errors.New("material not found")
	ErrSupplierAlreadyExist = errors.New("material already exist")
)

type Supplier struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 名称
	ShortName string `orm:"column(short_name)" json:"shortName"`
	// 长名称
	LongName string `orm:"column(long_name)" json:"longName"`
	// 英文名称
	EnglishName string `orm:"column(english_name)" json:"englishName"`
	// 法人
	LegalPerson string `orm:"column(legal_person)" json:"legalPerson"`
	// 企业代码
	EnterpriseCode string `orm:"column(enterprise_code)" json:"enterpriseCode"`
	// 联系人
	Contactor string `orm:"column(contactor)" json:"contactor"`
	// 联系电话
	Mobile string `orm:"column(mobile)" json:"mobile"`
	// 企业电话
	Telephone string `orm:"column(telephone)" json:"telephone"`
	// 邮编
	Postcode string `orm:"column(postcode)" json:"postcode"`
	// 传真
	Fax string `orm:"column(fax)" json:"fax"`
	// 邮箱
	Email string `orm:"column(email)" json:"email"`
	// 网址
	Website string `orm:"column(website)" json:"website"`
	// 省份
	Province string `orm:"column(province)" json:"province"`
	// 城市
	City string `orm:"column(city)" json:"city"`
	// 地区
	District string `orm:"column(district)" json:"district"`
	// 详细地址
	Address string `orm:"column(address)" json:"address"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *Supplier) TableName() string {
	return "supplier"
}

// 添加
func InsertSupplier(obj *Supplier) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelSupplier(id int) error {
	o := orm.NewOrm()

	obj := &Supplier{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateSupplier(obj *Supplier) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func SupplierById(id int) (*Supplier, error) {
	o := orm.NewOrm()

	obj := &Supplier{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrSupplierNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据名称查询
func SupplierByName(name string) (*Supplier, error) {
	o := orm.NewOrm()

	obj := &Supplier{
		ShortName: name,
	}

	if err := o.Read(obj, "ShortName"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrSupplierNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func SupplierList(where map[string]interface{}, page, pageSize int) (int64, []*Supplier, error) {
	o := orm.NewOrm()

	list := []*Supplier{}

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

		name := where["name"]
		if name != "" {
			sql += " AND (t1.short_name LIKE '%" + fmt.Sprintf("%s", name) + "%'"
			sql += " OR t1.long_name LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(supplierListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(supplierListSql+sql+" GROUP BY t1.id ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const supplierListCountSql = `
SELECT
    COUNT(*)
FROM
    supplier AS t1
WHERE
`

const supplierListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.short_name,
    t1.long_name,
    t1.english_name,
    t1.legal_person,
    t1.enterprise_code,
    t1.contactor,
    t1.mobile,
    t1.telephone,
    t1.postcode,
    t1.fax,
    t1.email,
    t1.website,
    t1.province,
    t1.city,
    t1.district,
    t1.address,
    t1.status,
    t1.updated,
    t1.updated_by
FROM
    supplier AS t1
WHERE
`
