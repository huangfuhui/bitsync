package account

import (
	"bitsync/controllers"
	"bitsync/validator"
	"bitsync/validator/account"
)

type AccountController struct {
	controllers.BaseController
}

func (c *AccountController) Register() {
	handset := c.GetString("handset")
	password := c.GetString("password")
	pin := c.GetString("pin")

	v := validator.BaseValidator{}
	v.Validate(&c.BaseController, account.Register{
		handset,
		password,
		pin,
	})

	c.Output("hello")
}
