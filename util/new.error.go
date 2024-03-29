package util

import "errors"

var (
	UsernameAlreadyExists      = errors.New("Username already exists")
	EmailAlreadyExists         = errors.New("Email already exists")
	UsernameOrEmailNotFound    = errors.New("Username or Email not found")
	InvalidPassword            = errors.New("Invalid password")
	EmailNotVerified           = errors.New("Email not verified")
	EmailVerifyAlreadyUsed     = errors.New("Email verify already used")
	EmailVerifyExpired         = errors.New("Email verify expired")
	EmailVerifyCodeNotValid    = errors.New("Email verify code not valid")
	UsernameNotFound           = errors.New("Username not found")
	UnauthorizedDeleteUser     = errors.New("Unauthorized delete user")
	AccountAlreadyExists       = errors.New("Account already exists")
	UserNotFound               = errors.New("User not found")
	UsernameNotExist           = errors.New("Username not exist")
	EmailNotExists             = errors.New("Email not exists")
	UsernameAndEmailNotMatch   = errors.New("Username and email not match")
	InvalidResetToken          = errors.New("Invalid secret token")
	ResetTokenAlreadyUsed      = errors.New("Secret token already used")
	AccountNotBelongToUser     = errors.New("Account not belong to user")
	InsufficientBalance        = errors.New("Insufficient balance")
	InvalidPin                 = errors.New("Invalid pin")
	DestinationAccountNotExist = errors.New("Destination account not exist")
	AccounDoesNotExist         = errors.New("Account does not exist")
	DateFormatNotValid         = errors.New("Date format not valid")
	CannotExportEmptyData      = errors.New("Cannot export empty data to PDF")
	FailedDeleteUserAccount    = errors.New("Failed delete user, account balance is not zero")
	ExpensesPlanNotFound       = errors.New("Expenses plan not found")
	SessionNotFound            = errors.New("Session not found")
	SessionExpired             = errors.New("Session expired")
	SessionIsBlocked           = errors.New("Session is blocked")
	SessionNotMatchUser        = errors.New("Session not match user")
	InvalidRefreshToken        = errors.New("Invalid refresh token")
)
