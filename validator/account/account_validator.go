package account

type Register struct {
	Handset  string `valid:"Required;Mobile"`
	Password string `valid:"Required"`
	Pin      string `valid:"Required;Numeric;Length(4)"`
}

type RegisterPIN struct {
	Handset string `valid:"Required;Mobile"`
}

type Login struct {
	Handset  string `valid:"Required;Mobile"`
	Password string `valid:"Required"`
}

type ModifyPassword struct {
	OldPwd string `valid:"Required"`
	NewPwd string `valid:"Required"`
}

type ResetPassword struct {
	Handset string `valid:"Required;Mobile"`
	Pin     string `valid:"Required;Numeric;Length(4)"`
}
