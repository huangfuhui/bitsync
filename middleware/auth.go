package middleware

type Auth struct {
	base Base
}

// 校验token有效性
func (auth *Auth) Verify() bool {

	return false
}
