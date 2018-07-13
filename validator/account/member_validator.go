package account

type Update struct {
	Name     string `valid:"Required;AlphaDash"`
	Email    string `valid:"Email"`
	Sex      int    `valid:"Numeric;Range(0, 1)"`
	Birthday string `valid:"Date"`
}
