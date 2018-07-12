package market

import (
	"bitsync/controllers"
	"bitsync/validator"
	"bitsync/validator/market"
	marketLogic "bitsync/logic/market"
)

type PriceController struct {
	controllers.BaseController
}

func (c *PriceController) Get() {
	exchange, _ := c.GetInt("exchange")
	symbolType, _ := c.GetInt("symbol_type")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, market.GetPrice{
		Exchange:   exchange,
		SymbolType: symbolType,
	})
	if !ok {
		return
	}

	l := marketLogic.PriceLogic{}
	l.Get(exchange, symbolType)

	c.Output("")
}
