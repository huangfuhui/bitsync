package sms

type AddTask struct {
	Type       int     `valid:"Required;Numeric;Range(1,1)"`
	ExchangeId int     `valid:"Required;Numeric;Range(1,2)"`
	SymbolPair string  `valid:"Required;Alpha"`
	Deviation  int     `valid:"Required;Numeric;Range(1,2)"`
	Value      float64 `valid:"Required;Numeric"`
}

type CancelTask struct {
	TaskId int `valid:"Required;Numeric;Min(1)"`
}
