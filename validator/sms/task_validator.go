package sms

type AddTask struct {
	Type       int    `valid:"Required;Range(1,1)"`
	ExchangeId int    `valid:"Required;Range(1,2)"`
	SymbolPair string `valid:"Required;AlphaDash"`
	Deviation  int    `valid:"Required;Range(1,2)"`
	Value      string `valid:"Required;Num"`
}

type CancelTask struct {
	TaskId int `valid:"Required;Min(1)"`
}

type GetTask struct {
	Type       int    `valid:"Required;Range(1,1)"`
	ExchangeId int    `valid:"Required;Range(1,2)"`
	SymbolPair string `valid:"Required;AlphaDash"`
}
