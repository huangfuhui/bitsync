package task

import (
	"github.com/astaxie/beego"
	"bitsync/models/sms"
	"strings"
	"bitsync/util"
	smsObj "bitsync/object/sms"
	"strconv"
	"github.com/astaxie/beego/orm"
	"bitsync/object"
	"time"
	"bitsync/services"
)

var exchange = map[int]string{
	1: "huobi",
	2: "dragonex",
	3: "okex",
	4: "binance",
	5: "gate",
	6: "bithumb",
}

var redis util.RedisCli

type WarnTask struct {
}

// 价格任务监控
func (task *WarnTask) Warn() {
	db, _ := beego.AppConfig.Int("redis_db_price_warn")
	var err error
	redis, err = util.NewPool(db)
	if err != nil {
		beego.Error("【Redis】", err)
		return
	}

	defer func() {
		if redis.Pool != nil {
			redis.ClosePool()
		}
	}()

	// 1.获取所有待执行任务
	m := sms.TaskThresholdValueModel{}
	taskList, err := m.WaitExecuteTaskList()
	if err != nil {
		beego.Error("【提醒任务】", err)
		return
	} else if len(taskList) == 0 {
		return
	}

	// 2.存放至Redis中
	for _, v := range taskList {
		id := v["id"]
		exchangeId, _ := strconv.Atoi(v["exchange_id"].(string))
		thresholdValue := v["threshold_value"]
		deviation := v["deviation"]

		key := strings.Replace(v["symbol_pair"].(string), "_", "", -1)
		value := exchange[exchangeId] + ":" + id.(string) + ":" + deviation.(string) + ":" + thresholdValue.(string)
		err := redis.SAdd(redis.Con(), key, value)
		if err != nil {
			beego.Error("【提醒任务】", err)
			return
		}

		err = redis.SAdd(redis.Con(), "task:threshold_value", key)
		if err != nil {
			beego.Error("【提醒任务】", err)
			return
		}
	}

	beego.Info("【提醒任务】初始化任务成功, 开始监控任务价格...")

	// 3.执行价格监控
	for {
		symbolPairs, err := redis.SMembers(redis.Con(), "task:threshold_value")
		if err != nil {
			beego.Error("【提醒任务】", err)
			return
		}

		for _, v := range symbolPairs {
			taskList, err := redis.SMembers(redis.Con(), v)
			if err != nil {
				beego.Error("【提醒任务】", err)
				return
			}
			if len(taskList) == 0 {
				redis.SRem(redis.Con(), "task:threshold_value", v)
			}

			for _, task := range taskList {
				info := strings.Split(task, ":")
				currentPrice, _ := util.Redis.Get(util.Redis.Con(), info[0]+":"+v)
				if currentPrice == "" {
					continue
				}

				price, _ := strconv.ParseFloat(currentPrice, 64)
				deviation, _ := strconv.Atoi(info[2])
				thresholdValue, _ := strconv.ParseFloat(info[3], 64)
				if (deviation == smsObj.DEVIATION_GT && price > 0 && thresholdValue >= price) || (deviation == smsObj.DEVIATION_LT && price > 0 && thresholdValue <= price) {
					beego.Info("【提醒任务】", task)
					go remind(task)
					redis.SRem(redis.Con(), v, task)
				}
			}
		}
	}
}

// 发短信提醒用户
func remind(taskInfo string) {
	tplId, _ := beego.AppConfig.Int64("sms::tpl_price_warn")
	info := strings.Split(taskInfo, ":")
	id, _ := strconv.Atoi(info[1])

	o := orm.NewOrm()
	o.Begin()

	// 1.查询任务
	query := `
select a.uid, a.id, a.task_id, a.type, a.status, b.symbol_pair, b.threshold_value, b.deviation, c.handset
from sms_task a
join task_threshold_value b
on a.task_id = b.id
join member c
on a.uid = c.uid
where a.id = ?
`
	var tasks []orm.Params
	_, err := o.Raw(query, id).Values(&tasks)
	if err != nil {
		o.Rollback()
		beego.Error("【提醒任务】", err)
		return
	} else if len(tasks) != 1 {
		o.Rollback()
		beego.Error("【提醒任务】查询任务错误: ", tasks)
		return
	}
	task := tasks[0]
	UID, _ := strconv.Atoi(task["uid"].(string))
	smsTaskId, _ := strconv.Atoi(task["id"].(string))
	symbolPair := task["symbol_pair"].(string)

	// 2.检查任务状态
	taskStatus, _ := strconv.Atoi(task["status"].(string))
	if taskStatus != smsObj.STATUS_WAIT {
		o.Rollback()

		beego.Debug("【提醒任务】任务非等待执行状态: ", task)
		redis.SRem(redis.Con(), strings.Replace(symbolPair, "_", "", -1), taskInfo)

		return
	}

	// 3.更新任务状态
	smsTask := smsObj.SmsTask{
		Base:       object.Base{Id: id},
		Status:     smsObj.STATUS_SUCCESS,
		FinishTime: time.Now(),
	}
	_, err = o.Update(&smsTask, "Status", "FinishTime")
	if err != nil {
		o.Rollback()
		beego.Error("【提醒任务】更新任务状态失败: ", err)
		return
	}

	// 4.扣减预消费
	query = `
update sms_wallet
set prepare_consume = prepare_consume - 1
where uid = ?
`
	_, err = o.Raw(query, UID).Exec()
	if err != nil {
		o.Rollback()
		beego.Error("【提醒任务】扣减预消费金额失败: ", err)
		return
	}

	err = o.Commit()
	if err != nil {
		beego.Error("【提醒任务】", err)
	}

	// 5.发送短信
	s := services.SmsService{}
	msgErr := s.SendSingle("86", task["handset"].(string)+"asd", []string{strings.Replace(symbolPair, "_", "/", -1), task["threshold_value"].(string)}, tplId)
	if msgErr != nil {
		beego.Error("【提醒任务】短信发送失败: ", msgErr)

		// 5.1退款
		err = refund(UID)
		if err != nil {
			beego.Error("【提醒任务】退款失败: ", err)
		}

		// 5.2标记任务失败
		smsTask = smsObj.SmsTask{
			Base:       object.Base{Id: id},
			Status:     smsObj.STATUS_FAIL,
			FinishTime: time.Now(),
		}
		_, err = o.Update(&smsTask, "Status", "FinishTime")
		if err != nil {
			beego.Error("【提醒任务】更新任务状态失败: ", err)
		}

		// 5.3记录失败任务
		failed := sms.SmsFailedTaskModel{}
		err = failed.Add(UID, smsTaskId, msgErr.Error())
		if err != nil {
			beego.Error("【提醒任务】", err)
		}
	}
}

// 退款
func refund(UID int) error {
	query := `
update sms_wallet
set balance = balance + 1
where uid = ?
`
	_, err := orm.NewOrm().Raw(query, UID).Exec()
	return err
}
