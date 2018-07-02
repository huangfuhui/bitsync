package account

import "time"

type Update struct {
	Name      string `valid:"Required;AlphaDash"`
	Email     string `valid:"Email"`
	Sex       int    `valid:"Numeric;Range(0, 1)"`
	AvatarUrl string `valid:"MaxSize(100)"`
	Birthday  time.Time
}
