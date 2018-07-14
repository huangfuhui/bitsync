package member

type Update struct {
	Name     string `valid:"MinSize(1)"`
	Email    string `valid:"Mail"`
	Sex      int    `valid:"Range(0,1)"`
	Birthday string `valid:"Date"`
}
