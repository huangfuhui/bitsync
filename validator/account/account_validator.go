package account

type Register struct {
	Handset  string `valid:"Required;Mobile"`
	Password string `valid:"Required"`
	Pin      string `valid:"Required;Numeric;Length(4)"`
}

type RegisterPIN struct {
	Handset  string `valid:"Required;Mobile"`
}
