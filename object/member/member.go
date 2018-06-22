package member

import (
	"bitsync/object"
	"time"
)

type Member struct {
	object.Base
	AccountId int
	Name      string
	HandSet   string
	Email     string
	Sex       int
	AvatarUrl string
	BirthDay  time.Time
}
