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
	l.Add(types, exchangeId, symbolPair, deviation, value)

	c.Output("")
}

func (c *TaskController) Get() {

}

func (c *TaskController) List() {

}

func (c *TaskController) Cancel() {

}
