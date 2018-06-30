package member

import (
	"bitsync/object"
	"time"
)

const (
	MEMBER_SEX_MALE   = 1
	MEMBER_SEX_FEMALE = 0
)

type Member struct {
	object.Base
	UID       int `orm:"column(UID)"`
	Name      string
	Handset   string
	Email     string
	Sex       int
	AvatarUrl string
	Birthday  time.Time
}
