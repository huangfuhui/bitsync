package sms

import (
	"bitsync/models"
	"bitsync/object/sms"
	"github.com/astaxie/beego/orm"
)

type SmsTaskModel struct {
	models.BaseModel
}

// 查询任务列表
func (m *SmsTaskModel) GetList(UID int) ([]orm.Params, error) {
	query := `
select a.type, c.name_cn name, b.symbol_pair, count(distinct a.id) task_quantity
from sms_task a
join task_threshold_value b
on a.task_id = b.id
and a.type = 1
join exchange c
on b.exchange_id = c.exchange_id
where a.status = 0
and a.uid = ?
group by b.exchange_id, b.symbol_pair
`

	var res []orm.Params
	_, err := orm.NewOrm().Raw(query, UID).Values(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// 查询某个任务
func (m *SmsTaskModel) Get(UID, taskId int) () {

}

// 查询任务状态
func (m *SmsTaskModel) Status(UID, taskId int) (status int, err error) {
	task := sms.SmsTask{
		UID:    UID,
		TaskId: taskId,
	}

	err = orm.NewOrm().Read(&task, "UID", "TaskId")
	if err != nil {
		return 0, err
	}

	return task.Status, nil
}

// 取消任务
func (m *SmsTaskModel) Cancel(UID, taskId int) error {
	o := orm.NewOrm()
	o.Begin()

	// 1.更新任务状态
	_, err := o.QueryTable("sms_task").
		Filter("uid", UID).
		Filter("id", taskId).
		Update(orm.Params{"status": sms.STATUS_CANCEL})
	if err != nil {
		o.Rollback()
		return err
	}

	// 2.退款
	query := `
update sms_wallet
set balance = balance +1, prepare_consume = prepare_consume - 1
where uid = ?
`
	_, err = o.Raw(query, UID).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	err = o.Commit()
	if err != nil {
		return err
	}

	return nil
}
