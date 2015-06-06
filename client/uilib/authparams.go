package uilib

const (
	AUTH_OP_LOGIN         = "login"
	AUTH_OP_LOGOUT        = "logout"
	AUTH_OP_PWD_RESET     = "pwdreset"
	AUTH_OP_PWD_RESET_REQ = "pwdresetreq"
)

//XXX UGH has to be copied from s5 *SERVER* library to our client side
type PasswordAuthParameters struct {
	Username         string
	Password         string
	ResetRequestUdid string
	UserUdid         string
	Op               string
}
