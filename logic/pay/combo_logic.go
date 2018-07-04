package pay

import (
	"bitsync/logic"
	"bitsync/models/pay"
	"github.com/astaxie/beego"
	payObj "bitsync/object/pay"
)

type ComboLogic struct {
	logic.BaseLogic
}

// 获取套餐
func (l *ComboLogic) Get() []payObj.Combo {
	comboModel := pay.ComboModel{}
	combos, err := comboModel.Get()
	if err != nil {
		beego.Error(err)
		l.Warn("获取套餐失败")
		return []payObj.Combo{}
	}

	return combos
}
