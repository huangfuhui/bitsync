package sms

import (
	"bitsync/models"
	"bitsync/object/sms"
	"github.com/astaxie/beego/orm"
)

type TaskThresholdValueModel struct {
	models.BaseModel
}

// 新建任务
func (m *TaskThresholdValueModel) Add(task *sms.TaskThresholdValue) (err error) {
	UID := task.UID

	o := orm.NewOrm()
	o.Begin()

	// 1.划扣钱包
	query := `
update sms_wallet
set balance = balance - 1, prepare_consume = prepare_consume + 1
where uid = ?
`
	_, err = o.Raw(query, UID).Exec()
	if err != nil {
		o.Rollback()

		return err
	}

	// 2.添加任务
	taskId, err := o.Insert(task)
	if err != nil {
		o.Rollback()

		return err
	}

	// 3.添加任务列表
	taskList := sms.SmsTask{
		UID:    UID,
		TaskId: int(taskId),
		Type:   sms.TYPE_THRESOLD_VALUE,
		Status: sms.STATUS_WAIT,
	}
	_, err = o.Insert(&taskList)
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

// 待执行的任务列表
func (m *TaskThresholdValueModel) WaitExecuteTaskList() ([]orm.Params, error) {
	query := `
select b.id, a.uid, a.symbol_pair, a.exchange_id, a.threshold_value, a.deviation
from task_threshold_value a
join sms_task b
on a.id = b.task_id
and b.status = 0
and b.type = 1
`
	var res []orm.Params
	_, err := orm.NewOrm().Raw(query).Values(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
