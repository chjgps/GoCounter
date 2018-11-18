package box

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrCorrectNotFound     = errors.New("grid not found")
	ErrCorrectAlreadyExist = errors.New("grid already exist")
)

type Correct struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 用户ID
	AccountId int `orm:"column(account_id)" json:"accountId"`
	// 格子ID
	GridId int `orm:"column(grid_id)" json:"gridId"`
	// 称重重量
	Weight int `orm:"column(weight)" json:"weight"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`

	// other
	GridName string `json:"gridName"`
}

func (t *Correct) TableName() string {
	return "rel_grid_correct"
}

// 添加
func InsertCorrect(obj *Correct) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 查询所有
func CorrectList(where map[string]interface{}, page, pageSize int) (int64, []*Correct, error) {
	o := orm.NewOrm()

	list := []*Correct{}

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

		accountId := where["accountId"]
		if accountId.(int) > 0 {
			sql += " AND t1.account_id = " + fmt.Sprintf("%d", accountId) + " "
		}

		gridId := where["gridId"]
		if gridId.(int) > 0 {
			sql += " AND t1.grid_id = " + fmt.Sprintf("%d", gridId) + " "
		}

		name := where["name"]
		if name != "" {
			sql += " AND (t1.name LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(correctListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(correctListSql+sql+" ORDER BY t1.id DESC LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const correctListCountSql = `
SELECT
    COUNT(*)
FROM
    rel_grid_correct AS t1
WHERE
`

const correctListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.account_id,
    t1.grid_id,
    t1.weight,
    t1.status,
    t1.updated,
    t1.updated_by,
    t2.name AS grid_name
FROM
    rel_grid_correct AS t1
LEFT JOIN
    rel_box_grid AS t2
ON
    t1.grid_id = t2.id
WHERE
`
