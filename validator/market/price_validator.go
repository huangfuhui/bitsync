package market

type GetPrice struct {
	Exchange int `valid:"Required;Numeric;Min(1)"`
	SymbolType int `valid:"Required;Numeric;Min(1)"`
}
