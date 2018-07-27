package task

import (
	"github.com/astaxie/beego"
	"bitsync/models/sms"
	"strings"
	"bitsync/util"
	smsObj "bitsync/object/sms"
	"strconv"
	"github.com/astaxie/beego/orm"
	"bitsync/services"
	"bitsync/object"
	"time"
)

var exchange = map[int]string{
	1: "huobi",
	2: "dragonex",
	3: "okex",
	4: "binance",
	5: "gate",
	6: "bithumb",
}

type WarnTask struct {
}

// 价格任务监控
func (task *WarnTask) Warn() {
	db, _ := beego.AppConfig.Int("redis_db_price_warn")
	redis := util.NewPool(db)
	con := redis.Con()

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
		exchangeId := v["exchange_id"]
		thresholdValue := v["threshold_value"]
		deviation := v["deviation"]

		key := strings.Replace(v["symbol_pair"].(string), "_", "", -1)
		value := exchange[exchangeId.(int)] + ":" + id.(string) + "：" + deviation.(string) + ":" + thresholdValue.(string)
		err := redis.SAdd(con, key, value)
		if err != nil {
			beego.Error("【提醒任务】", err)
			return
		}

		err = redis.SAdd(con, "task:threshold_value", key)
		if err != nil {
			beego.Error("【提醒任务】", err)
			return
		}
	}

	beego.Info("【提醒任务】初始化任务成功, 开始监控任务价格...")

	// 3.执行价格监控
	for {
		symbolPairs, err := redis.SMembers(con, "task:threshold_value")
		if err != nil {
			beego.Error("【提醒任务】", err)
			return
		}

		for _, v := range symbolPairs {
			taskList, err := redis.SMembers(con, v)
			if err != nil {
				beego.Error("【提醒任务】", err)
				return
			}

			for _, task := range taskList {
				info := strings.Split(task, ":")
				currentPrice, err := util.Redis.Get(util.Redis.Con(), info[0])
				if err != nil {
					beego.Error("【提醒任务】", err)
					return
				} else if currentPrice == "" {
					continue
				}

				price, _ := strconv.Atoi(currentPrice)
				deviation, _ := strconv.Atoi(info[2])
				thresholdValue, _ := strconv.Atoi(info[3])
				if (deviation == smsObj.DEVIATION_GT && price > 0 && thresholdValue >= price) || (deviation == smsObj.DEVIATION_LT && price > 0 && thresholdValue <= price) {
					go remind(task)
					redis.SRem(con, v, task)
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

	// 2.检查任务状态
	taskStatus, _ := strconv.Atoi(task["status"].(string))
	if taskStatus != smsObj.STATUS_WAIT {
		o.Rollback()
		beego.Error("【提醒任务】任务非等待执行状态: ", task)
		return
	}

	// 3.更新任务状态
	smsTask := smsObj.SmsTask{
		Base:       object.Base{Id: id},
		Status:     smsObj.STATUS_SUCCESS,
		FinishTime: time.Now(),
	}
	_, err = o.Update(smsTask, "Status", "FinishTime")
	if err != nil {
		o.Rollback()
		beego.Error("【提醒任务】更新任务状态失败: ", err)
		return
	}

	// 4.扣减预消费

	// 5.发送短信
	s := services.SmsService{}
	err = s.SendSingle("86", task["handset"].(string), []string{strings.Replace(task["symbol_pair"].(string), "_", "/", -1), task["threshold_value"].(string)}, tplId)
	if err != nil {
		o.Commit()
		beego.Error("【提醒任务】短信发送失败: ", err)

		// 退款

		return
	}
}
