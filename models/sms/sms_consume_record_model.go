package sms

import (
	"bitsync/models"
	"time"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type SmsConsumeRecordModel struct {
	models.BaseModel
}

func (m *SmsConsumeRecordModel) Get() {

}

func (m *SmsConsumeRecordModel) Add() {

}

// 按时间查询消费条数
func (m *SmsConsumeRecordModel) ConsumeQuantity(UID int, beginDate, endDate time.Time) (int, error) {
	query := `
select sum(amount) quantity
from sms_consume_record
where uid = ?
and status = 1
and date(consume_time) >= ?
and date(consume_time) <= ?
`

	var res []orm.Params
	_, err := orm.NewOrm().Raw(query, UID, beginDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Values(&res)
	if err != nil {
		return 0, err
	}

	if res[0]["quantity"] == nil {
		return 0, nil
	}

	quantity, err := strconv.Atoi(res[0]["quantity"].(string))
	if err != nil {
		return 0, err
	}

	return quantity, nil
}
