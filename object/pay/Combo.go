package pay

import "bitsync/object"

const (
	COMMO_STATUS_NO  = 0
	COMBO_STAUTS_YES = 1
)

type Combo struct {
	object.Base
	Name        string
	Price       int
	SmsQuantity int
	Description string
	Status      int
}
