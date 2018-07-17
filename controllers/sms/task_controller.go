package sms

import (
	"bitsync/controllers"
	"bitsync/validator"
	"bitsync/validator/sms"
	smsLogic "bitsync/logic/sms"
	"bitsync/logic"
)

type TaskController struct {
	controllers.BaseController
}

// 添加预警任务
func (c *TaskController) Add() {
	types, _ := c.GetInt("type")
	exchangeId, _ := c.GetInt("exchange_id")
	symbolPair := c.GetString("symbol_pair")
	deviation, _ := c.GetInt("deviation")
	value := c.GetString("value")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, sms.AddTask{
		types,
		exchangeId,
		symbolPair,
		deviation,
		value,
	})
	if !ok {
		return
	}

	l := smsLogic.TaskLogic{logic.BaseLogic{c.BaseController}}
	res := l.Add(types, exchangeId, symbolPair, deviation, value)

	c.Output(res)
}

// 获取某一交易对任务
func (c *TaskController) Get() {
	types, _ := c.GetInt("type")
	exchangeId, _ := c.GetInt("exchange_id")
	symbolPair := c.GetString("symbol_pair")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, sms.GetTask{
		types,
		exchangeId,
		symbolPair,
	})
	if !ok {
		return
	}

	l := smsLogic.TaskLogic{logic.BaseLogic{c.BaseController}}
	res := l.Get(types, exchangeId, symbolPair)

	c.Output(res)
}

// 获取任务列表
func (c *TaskController) List() {
	l := smsLogic.TaskLogic{logic.BaseLogic{c.BaseController}}
	res := l.List()

	c.Output(res)
}

// 取消任务
func (c *TaskController) Cancel() {
	taskId, _ := c.GetInt("task_id")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, sms.CancelTask{
		taskId,
	})
	if !ok {
		return
	}

	l := smsLogic.TaskLogic{logic.BaseLogic{c.BaseController}}
	l.Cancel(taskId)

	c.Output("")
}
