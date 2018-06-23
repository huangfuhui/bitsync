package object

import "time"

type Base struct {
	Id        int       `orm:"pk;auto"`
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}
