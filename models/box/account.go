package box

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrAccountAlreadyExist = errors.New("account already exist")
)

type Account struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 用户ID
	AccountId int `orm:"column(account_id)" json:"accountId"`
	// 柜子ID
	BoxId int `orm:"column(box_id)" json:"boxId"`
	// 格子ID
	GridId int `orm:"column(grid_id)" json:"gridId"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`

	// other
	GridName string `json:"gridName"`
}

func (t *Account) TableName() string {
	return "rel_account_grid"
}

// 添加
func InsertAccount(obj *Account) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelAccount(id int) error {
	o := orm.NewOrm()

	obj := &Account{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateAccount(obj *Account) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func AccountById(id int) (*Account, error) {
	o := orm.NewOrm()

	obj := &Account{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrAccountNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func AccountList(where map[string]interface{}, page, pageSize int) (int64, []*Grid, error) {
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

		accountIdStr := where["accountId"]
		if accountIdStr.(int) > 0 {
			accountId := fmt.Sprintf("%d", accountIdStr)
			sql += " AND t2.account_id = " + accountId + "  OR (t1.type = 1 AND (t2.account_id = " + accountId + " OR t2.account_id = '' OR t2.account_id ISNULL)) "
		}

		boxId := where["boxId"]
		if boxId.(int) > 0 {
			sql += " AND t1.box_id = " + fmt.Sprintf("%d", boxId) + " "
		}

		name := where["name"]
		if name != "" {
			sql += " AND (t1.name LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(accountListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(accountListSql+sql+" GROUP BY t3.grid_id ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const accountListCountSql = `
SELECT
    COUNT(*)
FROM
    rel_box_grid AS t1
LEFT JOIN
    rel_account_grid AS t2
ON
    t1.id = t2.grid_id
WHERE
`

const accountListSql = `
SELECT
    t1.id,
    t1.created,
    t1.created_by,
    t1.name,
    t1.box_id,
    t1.channel,
    COUNT(t3.id) AS qty,
    t1.status,
    t1.updated,
    t1.updated_by,
    t1.code,
    t1.material_id,
    t1.type,
    t1.safe_qty,
    t2.id AS account_grid_id,
    t2.account_id,
    COUNT(t3.id) AS total_qty
FROM
    rel_box_grid AS t1
LEFT JOIN
    rel_account_grid AS t2
ON
    t1.id = t2.grid_id
LEFT JOIN
    stock_detail AS t3
ON
    t1.id = t3.grid_id
WHERE
`
