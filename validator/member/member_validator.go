package member

type Update struct {
	Name     string `valid:"Required;AlphaDash"`
	Email    string `valid:"Email"`
	Sex      int    `valid:"Range(0,1)"`
	Birthday string `valid:"Date"`
}
